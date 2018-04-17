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
		Handler: &ConnectionHandler{Router: &RequestLineRouter{}},
	}
}

func MakeTCPServer(host string, port uint16) *TCPServer {
	return &TCPServer{
		Host:    host,
		Port:    port,
		Handler: &ConnectionHandler{Router: &RequestLineRouter{}},
	}
}

func MakeTCPServerWithHandler(host string, port uint16, handler Handler) *TCPServer {
	return &TCPServer{
		Host:    host,
		Port:    port,
		Handler: handler,
	}
}

type TCPServer struct {
	Host     string
	Port     uint16
	Handler  Handler
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
	server.Handler.Handle(reader, conn)
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

type Handler interface {
	Handle(requestReader *bufio.Reader, responseWriter io.Writer)
}

type Router interface {
	AddRoute(route Route)
	ParseRequest(reader *bufio.Reader) (ok Request, routeError Response)
}

type Request interface {
	Handle(client io.Writer) error
}

type Response interface {
	WriteTo(client io.Writer) error
}
