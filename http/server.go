package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Builder for TCPServer that defaults to any available port on localhost
func TCPServerBuilder(host string) *tcpServerBuilder {
	return &tcpServerBuilder{
		host:    host,
		port:    0,
		handler: NewConnectionHandler(NewRouter()),
	}
}

type tcpServerBuilder struct {
	host    string
	port    uint16
	handler ConnectionHandler
}

func (builder *tcpServerBuilder) Build() *TCPServer {
	return &TCPServer{
		Host:    builder.host,
		Port:    builder.port,
		Handler: builder.handler,
	}
}

func (builder *tcpServerBuilder) ListeningOnHost(host string) *tcpServerBuilder {
	builder.host = host
	return builder
}

func (builder *tcpServerBuilder) ListeningOnPort(port uint16) *tcpServerBuilder {
	builder.port = port
	return builder
}

func (builder *tcpServerBuilder) WithConnectionHandler(handler ConnectionHandler) *tcpServerBuilder {
	builder.handler = handler
	return builder
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

func (server *TCPServer) Routes() []Route {
	return server.Handler.Routes()
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

	//fmt.Printf("Listening for connections on %v\n", listener.Addr())
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

		go func() {
			server.Handler.Handle(bufio.NewReader(conn), conn)
			_ = conn.Close()
		}()
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
	Routes() []Route
}
