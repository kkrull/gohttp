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
			handler  *http.ConnectionHandler
			response = &mock.Response{}
			request  = &mock.Request{RespondReturns: response}
			router   = &mock.Router{ReturnsRequest: request}

			requestReader  = anyReader()
			responseWriter = anyWriter()
		)

		It("parses the request with the Router", func() {
			handler = &http.ConnectionHandler{Router: router}
			handler.NewHandle(requestReader, responseWriter)
			router.VerifyReceived(requestReader)
		})

		Context("when there is a routing error", func() {
			It("writes the routing error response to the response writer", func() {
				errorResponse := &mock.Response{}
				router = &mock.Router{ReturnsError: errorResponse}

				handler = &http.ConnectionHandler{Router: router}
				handler.NewHandle(requestReader, responseWriter)
				errorResponse.VerifyWrittenTo(responseWriter)
			})
		})

		It("responds to the request on the given writer", func() {
			handler = &http.ConnectionHandler{Router: router}
			handler.NewHandle(anyReader(), anyWriter())
			request.VerifyRespond()
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
