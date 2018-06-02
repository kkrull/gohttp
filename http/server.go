package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
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
	allRequestsHaveTaken := time.Duration(0)
	for i := 1; ; i++ {
		conn, listenerClosed := server.listener.AcceptTCP()
		if listenerClosed != nil {
			return
		}

		//connectionNumber := i
		startTime := time.Now()
		//fmt.Printf("Accept [%04d]: %v --> %v\n", connectionNumber, conn.RemoteAddr(), conn.LocalAddr())
		go func() {
			server.Handler.Handle(bufio.NewReader(conn), conn)
			endTime := time.Now()
			_ = conn.Close()
			thisRequestTook := endTime.Sub(startTime)
			allRequestsHaveTaken += thisRequestTook
			//fmt.Printf("Closed [%04d]: %v --> %v (%f ms / %f sec total)\n",
			//	connectionNumber,
			//	conn.RemoteAddr(),
			//	conn.LocalAddr(),
			//	float64(thisRequestTook.Nanoseconds())/float64(time.Millisecond),
			//	float64(allRequestsHaveTaken.Nanoseconds())/float64(time.Second),
			//)
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
