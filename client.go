package fwd

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type ForwardingConfiguration struct {
	LocalPort uint16
	Target    string
}

func RunClient(serverAddr string, interface_ string, forwards []ForwardingConfiguration) error {
	upstreamConn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer upstreamConn.Close()
	client := NewProxyClient(upstreamConn)

	var wg sync.WaitGroup
	for _, forward := range forwards {
		address := fmt.Sprintf("%s:%d", interface_, forward.LocalPort)
		wg.Add(1)
		go func() {
			runSingleClient(address, forward.Target, client)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
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
		log.Println(err)
		return
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
