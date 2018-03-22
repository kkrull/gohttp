package main

import (
	"fmt"
	"flag"
	"os"
	"net"
)

func main() {
	builder := ParseFromCommandLine()
	httpd, configError := builder.Build()
	if configError != nil {
		fmt.Println(configError.Error())
		os.Exit(1)
	}

	fmt.Printf("Starting httpd...\n")
	fmt.Printf("- Port: %d\n", httpd.Port)
	fmt.Printf("- Content Root: %s\n", httpd.ContentRootDirectory)
	httpd.Listen()
}

func ParseFromCommandLine() HttpServerBuilder {
	port := flag.Uint("p", 0, "The TCP port on which to listen")
	contentRootDirectory := flag.String("d", "", "The root content directory, from which to operate")

	flag.Parse()
	return HttpServerBuilder{
		ContentRootDirectory: *contentRootDirectory,
		Port:                 *port,
	}
}

type HttpServerBuilder struct {
	ContentRootDirectory string
	Port                 uint
}

func (builder HttpServerBuilder) Build() (httpd HttpServer, err error) {
	if builder.Port == 0 {
		return httpd, fmt.Errorf("gohttp: missing port")
	} else if builder.ContentRootDirectory == "" {
		return httpd, fmt.Errorf("gohttp: missing content root directory")
	}

	return HttpServer{
		ContentRootDirectory: builder.ContentRootDirectory,
		Port:                 builder.Port,
	}, nil
}

type HttpServer struct {
	ContentRootDirectory string
	Port                 uint
}

func (httpd HttpServer) Listen() error {
	address := fmt.Sprintf(":%d", httpd.Port)
	listener, listenError := net.Listen("tcp", address)
	if listenError != nil {
		return listenError
	}

	for { //TODO KDK: Monitor a given channel to request shutdown and exit gracefully
		conn, connectionError := listener.Accept()
		if connectionError != nil {
			fmt.Println(connectionError.Error()) //TODO KDK: Push to a buffered channel
		}

		fmt.Printf("Connected! %v\n", conn)
	}
}
