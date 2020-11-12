package fwd

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer(address string) {
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
