package http_test

import (
	"bufio"
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("blockingConnectionHandler", func() {
	Describe("#Handle", func() {
		var (
			handler http.ConnectionHandler
			request *mock.Request
			router  *mock.Router

			requestReader  = anyReader()
			responseWriter = anyWriter()
		)

		It("parses the request with the Router", func() {
			request = &mock.Request{}
			router = &mock.Router{ReturnsRequest: request}

			handler = http.NewConnectionHandler(router)
			handler.Handle(requestReader, responseWriter)
			router.VerifyReceived(requestReader)
		})

		Context("when there is a routing error", func() {
			It("writes the routing error response to the response writer", func() {
				errorResponse := &mock.Response{}
				router = &mock.Router{ReturnsError: errorResponse}

				handler = http.NewConnectionHandler(router)
				handler.Handle(requestReader, responseWriter)
				errorResponse.VerifyWrittenTo(responseWriter)
			})
		})

		It("handles the request", func() {
			request = &mock.Request{}
			router = &mock.Router{ReturnsRequest: request}

			handler = http.NewConnectionHandler(router)
			handler.Handle(requestReader, responseWriter)
			request.VerifyHandle(responseWriter)
		})

		Context("when there is an error handling the request", func() {
			It("responds with InternalServerError", func() {
				request = &mock.Request{HandleReturns: "bang"}
				router = &mock.Router{ReturnsRequest: request}

				handler = http.NewConnectionHandler(router)
				handler.Handle(requestReader, responseWriter)
				Expect(responseWriter.Buffered()).To(BeNumerically(">", 0))
			})
		})
	})
})

func anyReader() *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(""))
}

func anyWriter() *bufio.Writer {
	return bufio.NewWriter(bytes.NewBufferString(""))
}
