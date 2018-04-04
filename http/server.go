package http

import (
	"fmt"
	"net"
	"bufio"
	"io"
	"bytes"
)

func MakeTCPServerOnAvailablePort(contentRootDirectory string, host string) Server {
	return &TCPServer{
		Host: host,
		Port: 0,
	}
}

func MakeTCPServer(contentRootDirectory string, host string, port uint16) Server {
	return &TCPServer{
		Host: host,
		Port: port,
	}
}

/* TCPServer */

type TCPServer struct {
	Host     string
	Port     uint16
	listener *net.TCPListener
}

func (server *TCPServer) Address() net.Addr {
	if server.listener == nil {
		return nil
	}

	return server.listener.Addr()
}

func (server *TCPServer) Start() error {
	if err := server.startListening(); err != nil {
		return err
	}

	go server.acceptConnections()
	return nil
}

func (server *TCPServer) startListening() error {
	if server.listener != nil {
		return fmt.Errorf("TCPServer: already running")
	}

	address, addressErr := net.ResolveTCPAddr("tcp", server.hostAndPort())
	if addressErr != nil {
		return addressErr
	}

	listener, listenError := net.ListenTCP("tcp", address)
	if listenError != nil {
		return listenError
	}

	server.listener = listener
	return nil
}

func (server TCPServer) hostAndPort() string {
	return fmt.Sprintf("%s:%d", server.Host, server.Port)
}

func (server TCPServer) acceptConnections() {
	for {
		conn, listenerClosed := server.listener.AcceptTCP()
		if listenerClosed != nil {
			return
		}

		handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	_, err := readSocket(conn)
	if err != nil {
		return
	}

	fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n")
	_ = conn.Close()
}

func readSocket(conn *net.TCPConn) (*bytes.Buffer, error) {
	requestBytes := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	_, readError := reader.Read(requestBytes)
	if readError == io.EOF {
		return bytes.NewBuffer(nil), nil
	} else if readError != nil {
		return bytes.NewBuffer(nil), readError
	} else {
		return bytes.NewBuffer(requestBytes), nil
	}
}

func (server *TCPServer) Shutdown() error {
	if server.listener == nil {
		return nil
	}

	defer func() {
		server.listener = nil
	}()
	return server.listener.Close()
}

/* Server */

type Server interface {
	Address() net.Addr
	Start() error
	Shutdown() error
}
