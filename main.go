package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Command struct {
	Name string
	Code func()
}

var commands = []Command{
	{"proxy", mainProxy},
}

func printHelp(name string) {
	panic("TODO")
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
