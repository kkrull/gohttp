package http

import (
	"fmt"
	"net"
	"bufio"
)

func MakeTCPServerOnAvailablePort(contentRootDirectory string, host string) Server {
	return &TCPServer{
		Host:   host,
		Port:   0,
		Parser: RFC7230RequestParser{},
	}
}

func MakeTCPServer(contentRootDirectory string, host string, port uint16) Server {
	return &TCPServer{
		Host:   host,
		Port:   port,
		Parser: RFC7230RequestParser{},
	}
}

/* TCPServer */

type TCPServer struct {
	Host     string
	Port     uint16
	Parser   RequestParser
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

		server.handleConnection(conn)
		_ = conn.Close()
	}
}

func (server TCPServer) handleConnection(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	_, parseError := server.Parser.ParseRequest(reader)
	if parseError != nil {
		fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", parseError.StatusCode, parseError.Reason)
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
