package http_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "http")
}

/* HandlerMock */

type HandlerMock struct {
	handleRequestReader  *bufio.Reader
	handleResponseWriter io.Writer
}

func (mock *HandlerMock) Handle(requestReader *bufio.Reader, responseWriter io.Writer) {
	mock.handleRequestReader = requestReader
	mock.handleResponseWriter = responseWriter
}

func (mock *HandlerMock) ShouldHandleConnection() {
	ExpectWithOffset(1, mock.handleRequestReader).NotTo(BeNil())
	ExpectWithOffset(1, mock.handleResponseWriter).NotTo(BeNil())
}

func (mock *HandlerMock) Routes() []http.Route {
	return nil
}

/* PatchResourceMock */

type PatchResourceMock struct{}

func (mock *PatchResourceMock) Name() string {
	return "PatchResourceMock"
}

/* RequestLoggerMock */

type RequestLoggerMock struct {
	parsedReceived http.RequestMessage
}

func (mock *RequestLoggerMock) Parsed(message http.RequestMessage) {
	mock.parsedReceived = message
}

func (mock *RequestLoggerMock) ParsedShouldHaveReceived(method, target string) {
	ExpectWithOffset(1, mock.parsedReceived).NotTo(BeNil())
	ExpectWithOffset(1, mock.parsedReceived.Method()).To(Equal(method))
	ExpectWithOffset(1, mock.parsedReceived.Target()).To(Equal(target))
}

/* ResponseMock */

type ResponseMock struct {
	writeHeaderReceived io.Writer
	writtenTo           io.Writer
}

func (mock *ResponseMock) WriteTo(client io.Writer) error {
	mock.writtenTo = client
	return nil
}

func (mock *ResponseMock) WriteHeader(client io.Writer) error {
	mock.writeHeaderReceived = client
	return nil
}

func (mock *ResponseMock) VerifyWrittenTo(writer *bufio.Writer) {
	ExpectWithOffset(1, mock.writtenTo).To(BeIdenticalTo(writer))
}

/* RouteMock */

type RouteMock struct {
	RouteReturns   http.Request
	routeRequested http.RequestMessage
}

func (mock *RouteMock) Route(requested http.RequestMessage) http.Request {
	mock.routeRequested = requested
	return mock.RouteReturns
}

func (mock *RouteMock) ShouldHaveReceived(method string, target string) {
	ExpectWithOffset(1, mock.routeRequested.Method()).To(Equal(method))
	ExpectWithOffset(1, mock.routeRequested.Target()).To(Equal(target))
}

func (mock *RouteMock) ShouldHaveReceivedParameters(parameters map[string]string) {
	ExpectWithOffset(1, mock.routeRequested.QueryParameters()).To(Equal(parameters))
}

/* RouterMock */

type RouterMock struct {
	ReturnsRequest http.Request
	ReturnsError   http.Response
	receivedReader *bufio.Reader
	parsed         []byte
}

func (mock *RouterMock) Routes() []http.Route {
	return nil
}

func (mock *RouterMock) RouteRequest(reader *bufio.Reader) (http.Request, http.Response) {
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

func (mock RouterMock) VerifyReceived(reader *bufio.Reader) {
	ExpectWithOffset(1, mock.receivedReader).To(BeIdenticalTo(reader))
}

/* Helpers */

func makeReader(template string, values ...interface{}) *bufio.Reader {
	text := fmt.Sprintf(template, values...)
	return bufio.NewReader(bytes.NewBufferString(text))
}
