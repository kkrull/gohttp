package mock

import (
	"fmt"
	. "github.com/onsi/gomega"
	"net"
	"github.com/kkrull/gohttp/http"
	"bufio"
)

type RequestParser struct {
	ParseRequestReturnError *http.ParseError
	parseRequestReceived    []byte
}

func (parser *RequestParser) ParseRequest(reader *bufio.Reader) (*http.Request, *http.ParseError) {
	allButLF, _ := reader.ReadBytes(byte('\r'))
	shouldBeLF, _ := reader.ReadByte()
	parser.parseRequestReceived = appendByte(allButLF, shouldBeLF)
	return nil, parser.ParseRequestReturnError
}

func appendByte(allButLast []byte, last byte) []byte {
	whole := make([]byte, len(allButLast)+1)
	copy(whole, allButLast)
	whole[len(whole)-1] = last
	return whole
}

func (parser RequestParser) VerifyReceived(expected []byte) {
	Expect(parser.parseRequestReceived).To(Equal(expected))
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
