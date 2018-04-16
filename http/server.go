package http

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/kkrull/gohttp/msg/servererror"
)

func MakeTCPServerOnAvailablePort(contentRootDirectory string, host string) *TCPServer {
	return &TCPServer{
		Host:   host,
		Port:   0,
		Parser: RFC7230RequestParser{BaseDirectory: contentRootDirectory},
	}
}

func MakeTCPServer(contentRootDirectory string, host string, port uint16) *TCPServer {
	return &TCPServer{
		Host:   host,
		Port:   port,
		Parser: RFC7230RequestParser{BaseDirectory: contentRootDirectory},
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
	request, parseError := server.Parser.ParseRequest(reader)
	if parseError != nil {
		parseError.WriteTo(conn)
		return
	}

	requestError := request.Handle(conn)
	if requestError != nil {
		response := servererror.InternalServerError{}
		response.WriteTo(conn)
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

type RequestParser interface {
	ParseRequest(reader *bufio.Reader) (ok Request, parseError Response)
}

type Request interface {
	Handle(client io.Writer) error
}

type Response interface {
	WriteTo(client io.Writer) error
}
