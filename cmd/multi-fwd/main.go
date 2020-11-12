package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/edofic/kube-multi-fwd"
	"gopkg.in/yaml.v2"
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

type ClientConfiguration struct {
	Forwards []fwd.ForwardingConfiguration `yaml:"forwards"`
}

func mainClient() {
	serverF := flag.String("server", "127.0.0.1:50051", "Proxy server address")
	interfaceF := flag.String("interface", "127.0.0.1", "Interface to bind to")
	configF := flag.String("config", "", "Forwarding configuration file. Cannot be set if `forwards` is set")
	forwardsF := flag.String("forwards", "", "Comma separated list of forwards of the form <LOCAL PORT>:<TARGET HOST>:<TARGET PORT>")
	flag.Parse()

	if (*configF != "") && (*forwardsF != "") {
		log.Panic("Cannot specify both a config file and forwards flag")
	}

	var forwards []fwd.ForwardingConfiguration
	if *forwardsF != "" {
		forwards = parseForwards(*forwardsF)
	} else {
		config, err := ioutil.ReadFile(*configF)
		if err != nil {
			panic(err)
		}
		var configuration ClientConfiguration
		err = yaml.Unmarshal(config, &configuration)
		if err != nil {
			panic(err)
		}
		forwards = configuration.Forwards
	}

	if len(forwards) == 0 {
		log.Println("nothing to forward")
	}

	fwd.RunClient(*serverF, *interfaceF, forwards)
}

func parseForwards(raw string) []fwd.ForwardingConfiguration {
	var forwards []fwd.ForwardingConfiguration
	for _, forward := range strings.Split(raw, ",") {
		parts := strings.SplitN(forward, ":", 2)
		if len(parts) != 2 {
			log.Panic("Cannot parse port forward config", forward)
		}
		port, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}

		target := parts[1]
		forwards = append(forwards, fwd.ForwardingConfiguration{
			LocalPort: uint16(port),
			Target:    target,
		})
	}
	return forwards
}
