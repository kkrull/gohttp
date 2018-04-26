package mock

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/gomega"
)

type Handler struct {
	handleRequestReader  *bufio.Reader
	handleResponseWriter io.Writer
}

func (mock *Handler) Handle(requestReader *bufio.Reader, responseWriter io.Writer) {
	mock.handleRequestReader = requestReader
	mock.handleResponseWriter = responseWriter
}

func (mock *Handler) ShouldHandleConnection() {
	Expect(mock.handleRequestReader).NotTo(BeNil())
	Expect(mock.handleResponseWriter).NotTo(BeNil())
}

func (mock *Handler) Routes() []http.Route {
	return nil
}

type Router struct {
	ReturnsRequest http.Request
	ReturnsError   http.Response
	receivedReader *bufio.Reader
	parsed         []byte
}

func (mock *Router) Routes() []http.Route {
	return nil
}

func (mock *Router) RouteRequest(reader *bufio.Reader) (http.Request, http.Response) {
	mock.receivedReader = reader
	allButLF, _ := reader.ReadBytes(byte('\r'))
	shouldBeLF, _ := reader.ReadByte()
	mock.parsed = appendByte(allButLF, shouldBeLF)
	return mock.ReturnsRequest, mock.ReturnsError
}

func appendByte(allButLast []byte, last byte) []byte {
	whole := make([]byte, len(allButLast)+1)
	copy(whole, allButLast)
	whole[len(whole)-1] = last
	return whole
}

func (mock Router) VerifyReceived(reader *bufio.Reader) {
	Expect(mock.receivedReader).To(BeIdenticalTo(reader))
}

type Request struct {
	HandleReturns  string
	handleReceived io.Writer
}

func (mock *Request) Handle(writer io.Writer) error {
	mock.handleReceived = writer
	if mock.HandleReturns != "" {
		return fmt.Errorf(mock.HandleReturns)
	}

	return nil
}

func (mock *Request) VerifyHandle(writer *bufio.Writer) {
	ExpectWithOffset(1, mock.handleReceived).To(BeIdenticalTo(writer))
}

type Route struct {
	RouteReturns   http.Request
	routeRequested *http.RequestLine
}

func (mock *Route) Route(requested *http.RequestLine) http.Request {
	mock.routeRequested = requested
	return mock.RouteReturns
}

func (mock *Route) ShouldHaveReceived(method string, target string) {
	ExpectWithOffset(1, mock.routeRequested).To(BeEquivalentTo(&http.RequestLine{
		TheMethod: method,
		TheTarget: target,
	}))
}

func (mock *Route) ShouldHaveReceivedParameters(parameters map[string]string) {
	ExpectWithOffset(1, mock.routeRequested.TheQueryParameters).To(Equal(parameters))
}

type Response struct {
	writeHeaderReceived io.Writer
	writtenTo           io.Writer
}

func (mock *Response) WriteTo(client io.Writer) error {
	mock.writtenTo = client
	return nil
}

func (mock *Response) WriteHeader(client io.Writer) error {
	mock.writeHeaderReceived = client
	return nil
}

func (mock *Response) VerifyWrittenTo(writer *bufio.Writer) {
	ExpectWithOffset(1, mock.writtenTo).To(BeIdenticalTo(writer))
}
