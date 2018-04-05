package http

import (
	"fmt"
	"net"
	"bufio"
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
		_ = conn.Close()
	}
}

func handleConnection(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	methodBytes, _ := reader.ReadBytes(' ')
	if len(methodBytes) == 1 {
		fmt.Fprint(conn, "HTTP/1.1 400 Bad Request\r\n")
		return
	} else if len(methodBytes) > 8000 {
		fmt.Fprint(conn, "HTTP/1.1 501 Not Implemented\r\n")
		return
	}

	remainingBytes := make([]byte, 1024)
	_, _ = reader.Read(remainingBytes)
	fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n")
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
