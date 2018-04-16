package mock

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/gomega"
)

type Router struct {
	ReturnsRequest http.Request
	ReturnsError   http.Response
	received       []byte
}

func (mock *Router) AddRoute(route http.Route) {}

func (mock *Router) ParseRequest(reader *bufio.Reader) (http.Request, http.Response) {
	allButLF, _ := reader.ReadBytes(byte('\r'))
	shouldBeLF, _ := reader.ReadByte()
	mock.received = appendByte(allButLF, shouldBeLF)
	return mock.ReturnsRequest, mock.ReturnsError
}

func appendByte(allButLast []byte, last byte) []byte {
	whole := make([]byte, len(allButLast)+1)
	copy(whole, allButLast)
	whole[len(whole)-1] = last
	return whole
}

func (mock Router) VerifyReceived(expected []byte) {
	Expect(mock.received).To(Equal(expected))
}

type Request struct {
	ReturnsError string
}

func (mock Request) Handle(connWriter io.Writer) error {
	if mock.ReturnsError != "" {
		return fmt.Errorf(mock.ReturnsError)
	}

	return nil
}

type Route struct {
	routeMethod  string
	routeTarget  string
	RouteReturns http.Request
}

func (mock *Route) Route(method string, target string) http.Request {
	mock.routeMethod = method
	mock.routeTarget = target
	return mock.RouteReturns
}

func (mock *Route) ShouldHaveReceived(method string, target string) {
	Expect(mock.routeMethod).To(Equal(method))
	Expect(mock.routeTarget).To(Equal(target))
}

type Server struct {
	StartFails  string
	startCalled bool

	ShutdownFails  string
	shutdownCalled bool
}

func (Server) Address() net.Addr {
	panic("implement me")
}

func (mock *Server) Start() error {
	mock.startCalled = true
	if mock.StartFails != "" {
		return fmt.Errorf(mock.StartFails)
	}

	return nil
}

func (mock Server) VerifyStart() {
	Expect(mock.startCalled).To(BeTrue())
}

func (mock *Server) Shutdown() error {
	mock.shutdownCalled = true
	if mock.ShutdownFails != "" {
		return fmt.Errorf(mock.ShutdownFails)
	}

	return nil
}

func (mock Server) VerifyRunning() {
	Expect(mock.startCalled).To(BeTrue())
	Expect(mock.shutdownCalled).To(BeFalse())
}

func (mock Server) VerifyShutdown() {
	Expect(mock.shutdownCalled).To(BeTrue())
}
