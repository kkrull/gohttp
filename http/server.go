package http

import (
	"bufio"
	"fmt"
	"net"
)

func NewServer(contentRootDirectory string, port uint) Server {
	return server{ContentRootDirectory: contentRootDirectory, Port: port}
}

type Server interface {
	Listen() error
}

type server struct {
	ContentRootDirectory string
	Port                 uint
}

func (server server) Listen() error {
	address, addressErr := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", server.Port))
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
