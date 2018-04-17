package http_test

import (
	"bufio"
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/mock"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Handler", func() {
	Describe("#NewHandle", func() {
		var (
			handler *http.ConnectionHandler
			router  *mock.Router
		)

		It("parses the request with the Router", func() {
			router = &mock.Router{ReturnsRequest: &mock.Request{RespondReturns: &mock.Response{}}}
			requestReader := makeReader("any request")

			handler = &http.ConnectionHandler{Router: router}
			handler.NewHandle(requestReader, anyWriter())
			router.VerifyReceived(requestReader)
		})

		Context("when there is a routing error", func() {
			It("writes the routing error response to the response writer", func() {
				errorResponse := &mock.Response{}
				router = &mock.Router{ReturnsError: errorResponse}

				handler = &http.ConnectionHandler{Router: router}
				responseWriter := anyWriter()
				handler.NewHandle(anyReader(), responseWriter)
				errorResponse.VerifyWrittenTo(responseWriter)
			})
		})

		It("determines the response to the request", func() {
			request := &mock.Request{RespondReturns: &mock.Response{}}
			router = &mock.Router{ReturnsRequest: request}

			handler = &http.ConnectionHandler{Router: router}
			handler.NewHandle(anyReader(), anyWriter())
			request.VerifyRespond()
		})

		It("writes a response to a routed request", func() {
			response := &mock.Response{}
			request := &mock.Request{RespondReturns: response}
			router = &mock.Router{ReturnsRequest: request}

			handler = &http.ConnectionHandler{Router: router}
			responseWriter := anyWriter()
			handler.NewHandle(anyReader(), responseWriter)
			response.VerifyWrittenTo(responseWriter)
		})
	})
})

func anyReader() *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(""))
}

func anyWriter() *bufio.Writer {
	return bufio.NewWriter(bytes.NewBufferString(""))
}
