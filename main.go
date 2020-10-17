//go:generate protoc -I=. --go_out=. --go-grpc_out=. --go-grpc_opt=requireUnimplementedServers=false ./protocol.proto
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type Command struct {
	Name        string
	Description string
	Code        func()
}

var commands = []Command{
	{"proxy-test", "Quick proxying test", mainProxy},
	{"server", "Server side that proxies traffic", mainServer},
	{"client", "Client side that connects to the proxy", mainClient},
}

func printHelp(name string) {
	fmt.Printf("usage: %s COMMAND ...flags\n\n", name)
	fmt.Println("Available commands:")
	for _, command := range commands {
		fmt.Println(command.Name, strings.Repeat(" ", 10-len(command.Name)), command.Description)
	}
}

func main() {
	name := os.Args[0]
	if len(os.Args) < 2 {
		printHelp(name)
		return
	}
	selectedCommand := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)
	for _, command := range commands {
		if command.Name == selectedCommand {
			command.Code()
			return
		}
	}
	// command not found
	printHelp(name)
}

func mainProxy() {
	interfaceF := flag.String("interface", "127.0.0.1", "Interface to bind to")
	portF := flag.Int("port", 63000, "Port to listen on")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *interfaceF, *portF)
	log.Println("running proxy on", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go proxyConn(conn)
	}
}

func proxyConn(conn net.Conn) {
	target, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Println("Error dialing", err)
		return
	}
	closeCh := make(chan struct{}, 2)
	go pipeConn(conn, target, closeCh)
	go pipeConn(target, conn, closeCh)
	<-closeCh
	conn.Close()
	target.Close()
}

func pipeConn(source, target net.Conn, closeCh chan struct{}) {
	io.Copy(source, target)
	//log.Print(n, err)
	closeCh <- struct{}{}
}

func mainServer() {
	interfaceF := flag.String("interface", "127.0.0.1", "Interface to bind to")
	portF := flag.Int("port", 50051, "Port to listen on")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *interfaceF, *portF)
	log.Println("starting grpc server on", address)

	s := grpc.NewServer()
	RegisterProxyServer(s, NewProxy())
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Println(err)
	}
}

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Proxy(stream Proxy_ProxyServer) error {
	rawReq, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("received request", rawReq)
	switch req := rawReq.Req.(type) {
	case *ProxyRequest_Connect:
		log.Println("connecting", req)
		err = stream.Send(&ProxyResponse{Res: &ProxyResponse_Connected{}})
		if err != nil {
			log.Println("error sending connected response", err)
			return err
		}
	default:
		return errors.New("unknown request")
	}
	for {
		rawReq, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("received request", rawReq)
		switch req := rawReq.Req.(type) {
		case *ProxyRequest_Chunk:
			log.Println("echoing chunk", req)
			err = stream.Send(
				&ProxyResponse{
					Res: &ProxyResponse_Chunk{
						Chunk: req.Chunk,
					},
				},
			)
			if err != nil {
				log.Println(err)
				return err
			}
		default:
			return errors.New("unknown request")
		}
	}
}

func mainClient() {
	serverF := flag.String("server", "127.0.0.1:50051", "Proxy server address")
	interfaceF := flag.String("interface", "127.0.0.1", "Interface to bind to")
	portF := flag.Int("port", 63000, "Port to listen on")
	flag.Parse()

	upstreamConn, err := grpc.Dial(*serverF, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer upstreamConn.Close()
	client := NewProxyClient(upstreamConn)

	address := fmt.Sprintf("%s:%d", *interfaceF, *portF)
	log.Println("running proxy on", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go proxyConnOverGrpc(conn, client)
	}

}

func proxyConnOverGrpc(conn net.Conn, client ProxyClient) {
	defer conn.Close()
	proxyClient, err := client.Proxy(context.Background())
	if err != nil {
		log.Panic(err)
	}
	defer proxyClient.CloseSend()

	err = proxyClient.Send(&ProxyRequest{
		Req: &ProxyRequest_Connect{},
	})
	if err != nil {
		log.Panic(err)
	}

	resp, err := proxyClient.Recv()
	if err != nil {
		log.Panic(err)
	}
	if _, ok := resp.GetRes().(*ProxyResponse_Connected); !ok {
		log.Println("Did not connect")
		return
	}
	log.Println("connected")
	go func() {
		for {
			log.Println("reading chunks")
			resp, err := proxyClient.Recv()
			if err != nil {
				log.Println("error receiving", err)
				return
			}
			chunk, ok := resp.Res.(*ProxyResponse_Chunk)
			if !ok {
				log.Println("Failed to read chunk", resp.Res)
				return
			}
			buf := chunk.Chunk.Payload
			_, err = conn.Write(buf)
			if err != nil {
				log.Println("Failed to write to conn", err)
				conn.Close()
				return
			}
		}
	}()
	buf := make([]byte, 32*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		proxyClient.Send(
			&ProxyRequest{
				Req: &ProxyRequest_Chunk{
					Chunk: &ProxyBytes{
						Payload: buf[:n],
					},
				},
			},
		)
	}
}
