package main

import (
	"fmt"
	"flag"
	"os"
	"net"
	"bufio"
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
	address, addressErr := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", httpd.Port))
	if addressErr != nil {
		return addressErr
	}

	listener, listenError := net.ListenTCP("tcp", address)
	if listenError != nil {
		return listenError
	}

	for {
		conn, connectionError := listener.AcceptTCP()
		if connectionError != nil {
			fmt.Println(connectionError.Error())
		}

		fmt.Printf("Connected! %v -> %v\n", conn.LocalAddr(), conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	//Read so it doesn't complain about the connection being reset
	readBuffer := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	numBytesRead, readError := reader.Read(readBuffer)
	if readError != nil {
		fmt.Printf(readError.Error())
		return
	}

	fmt.Printf("Read %d bytes\n", numBytesRead)

	//Try first to respond with just enough for a 404
	fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n")
	if closeError := conn.Close(); closeError != nil {
		fmt.Printf("Failed to close connection: %s\n", closeError.Error())
	}
}
