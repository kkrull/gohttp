package http

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/kkrull/gohttp/msg/servererror"
)

func MakeTCPServerOnAvailablePort(host string) *TCPServer {
	return &TCPServer{
		Host:   host,
		Port:   0,
		Router: &RequestLineRouter{},
	}
}

func MakeTCPServer(host string, port uint16) *TCPServer {
	return &TCPServer{
		Host:   host,
		Port:   port,
		Router: &RequestLineRouter{},
	}
}

type TCPServer struct {
	Host     string
	Port     uint16
	Router   Router
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

func (server *TCPServer) AddRoute(route Route) {
	server.Router.AddRoute(route)
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
	request, parseError := server.Router.ParseRequest(reader)
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
