package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime"
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
	numTokens := 8
	tokens := make(chan int, numTokens)
	for token := 1; token <= numTokens; token++ {
		tokens <- token
	}

	//TODO KDK: Is there a problem with Apache Bench?
	//When things appear to hang, the server still responds to requests on all its tokens
	//Other times, it refuses connections, which suggests a problem with the code, the Go runtime, or the OS's network stack
	for c := 1; ; {
		select {
		case token := <-tokens:
			request := c
			c++
			conn, listenerClosed := server.listener.AcceptTCP()
			if listenerClosed != nil {
				fmt.Printf("Listener closed prematurely: %s\n", listenerClosed)
				return
			}

			go func(t int) {
				fmt.Printf("[%d/%04d] Handling request %d\n", t, runtime.NumGoroutine(), request)
				handleError := HardCodedResponse(conn)
				if handleError != nil {
					fmt.Printf("Error handling request: %s\n", handleError)
				}

				//server.Handler.Handle(bufio.NewReader(conn), conn)
				closeError := conn.Close()
				if closeError != nil {
					fmt.Printf("Close error: %s\n", closeError)
				}

				tokens <- t
			}(token)
		default:
			//runtime.Gosched()
			time.Sleep(5 * time.Millisecond)
			//fmt.Printf("Waited for token\n")
		}
	}
}

var (
	request = make([]byte, 1024)
)

func HardCodedResponse(conn *net.TCPConn) error {
	_, readError := conn.Read(request)
	//if numRead != 88 {
	//	fmt.Printf("Read %d bytes\n", numRead)
	//}
	if readError != nil {
		return readError
	}

	_, writeError := fmt.Fprintf(conn, "HTTP/1.1 204 No Content\r\n\r\n")
	//if numWritten != 27 {
	//	fmt.Printf("Wrote %d bytes\n", numWritten)
	//}
	if writeError != nil {
		return writeError
	}

	return nil
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
