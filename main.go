//go:generate protoc -I=. --go_out=. --go-grpc_out=. --go-grpc_opt=requireUnimplementedServers=false ./protocol.proto
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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
	var target net.Conn
	rawReq, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("received request", rawReq)
	switch req := rawReq.Req.(type) {
	case *ProxyRequest_Connect:
		log.Println("connecting", req)
		target, err = net.Dial("tcp", req.Connect.Target)
		if err != nil {
			log.Println("error connecting", err)
		}
		defer target.Close()
		err = stream.Send(&ProxyResponse{Res: &ProxyResponse_Connected{}})
		if err != nil {
			log.Println("error sending connected response", err)
			return err
		}
	default:
		return errors.New("unknown request")
	}
	go func() {
		defer target.Close()
		defer func() {
			stream.Send(&ProxyResponse{Res: &ProxyResponse_Eof{}})
		}()
		for {
			buf := make([]byte, 32*1024)
			n, err := target.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}
			chunk := buf[:n]
			err = stream.Send(
				&ProxyResponse{
					Res: &ProxyResponse_Chunk{
						Chunk: chunk,
					},
				},
			)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	for {
		rawReq, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		switch req := rawReq.Req.(type) {
		case *ProxyRequest_Chunk:
			_, err = target.Write(req.Chunk)
			if err != nil {
				log.Println(err)
				return err
			}
		default:
			log.Println("received request", rawReq)
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
	runSingleClient(address, "localhost:8000", client)
}

func runSingleClient(address, target string, client ProxyClient) {
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
		go proxyConnOverGrpc(target, conn, client)
	}
}

func proxyConnOverGrpc(target string, conn net.Conn, client ProxyClient) {
	defer conn.Close()
	proxyClient, err := client.Proxy(context.Background())
	if err != nil {
		log.Panic(err)
	}
	defer proxyClient.CloseSend()

	err = proxyClient.Send(&ProxyRequest{
		Req: &ProxyRequest_Connect{
			Connect: &ProxyConnect{Target: target},
		},
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
			resp, err := proxyClient.Recv()
			if err != nil {
				log.Println("error receiving", err)
				return
			}
			switch res := resp.Res.(type) {
			case *ProxyResponse_Chunk:
				_, err = conn.Write(res.Chunk)
				if err != nil {
					log.Println("Failed to write to conn", err)
					conn.Close()
					return
				}
			case *ProxyResponse_Eof:
				log.Println("Proxy EOF", res)
				return
			default:
				log.Println("Failed to read response", res)
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
					Chunk: buf[:n],
				},
			},
		)
	}
}
