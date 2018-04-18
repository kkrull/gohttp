package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func MakeTCPServerOnAvailablePort(host string) *TCPServer {
	return &TCPServer{
		Host:    host,
		Port:    0,
		Handler: &BlockingConnectionHandler{Router: &RequestLineRouter{}},
	}
}

func MakeTCPServer(host string, port uint16) *TCPServer {
	return &TCPServer{
		Host:    host,
		Port:    port,
		Handler: &BlockingConnectionHandler{Router: &RequestLineRouter{}},
	}
}

func MakeTCPServerWithHandler(host string, port uint16, handler ConnectionHandler) *TCPServer {
	return &TCPServer{
		Host:    host,
		Port:    port,
		Handler: handler,
	}
}

type TCPServer struct {
	Host     string
	Port     uint16
	Handler  ConnectionHandler
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

		server.Handler.Handle(bufio.NewReader(conn), conn)
		_ = conn.Close()
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

type ConnectionHandler interface {
	Handle(request *bufio.Reader, response io.Writer)
}
