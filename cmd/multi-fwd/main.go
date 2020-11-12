package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edofic/kube-multi-fwd"
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

	fwd.RunServer(address)
}

func mainClient() {
	serverF := flag.String("server", "127.0.0.1:50051", "Proxy server address")
	interfaceF := flag.String("interface", "127.0.0.1", "Interface to bind to")
	forwards := flag.String("forwards", "", "Comma separated list of forwards of the form <LOCAL PORT>:<TARGET HOST>:<TARGET PORT>")
	flag.Parse()

	fwd.RunClient(*serverF, *interfaceF, *forwards)
}
