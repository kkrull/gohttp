package http_test

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/kkrull/gohttp/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("RFC7230RequestParser", func() {
	Describe("#ParseRequest", func() {
		var (
			parser  *http.RFC7230RequestParser
			request http.Request
			err     *http.ParseError
		)

		BeforeEach(func() {
			parser = &http.RFC7230RequestParser{BaseDirectory: "/tmp"}
		})

		Describe("it returns 400 Bad Request", func() {
			It("for a completely blank request", func() {
				request, err = parser.ParseRequest(makeReader(""))
				Expect(err).To(beABadRequestError())
			})

			It("for any line missing CR", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n\n"))
				Expect(err).To(beABadRequestError())
			})

			It("for any line missing LF", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r"))
				Expect(err).To(beABadRequestError())
			})

			It("for a request missing a request-line", func() {
				request, err = parser.ParseRequest(makeReader("\r\n"))
				Expect(err).To(beABadRequestError())
			})

			It("for a request missing an ending CRLF", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n"))
				Expect(err).To(beABadRequestError())
			})

			It("when multiple spaces are separating fields in request-line", func() {
				request, err = parser.ParseRequest(makeReader("GET /  HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestError())
			})

			It("when fields in request-line contain spaces", func() {
				request, err = parser.ParseRequest(makeReader("GET /a\\ b HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestError())
			})

			It("when the request starts with whitespace", func() {
				request, err = parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestError())
			})
		})

		Context("given a well-formed GET request", func() {
			var (
				reader       *bufio.Reader
				typedRequest *http.GetRequest
			)

			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)

				parser = &http.RFC7230RequestParser{BaseDirectory: "/public"}
				request, err = parser.ParseRequest(reader)
				typedRequest = request.(*http.GetRequest)
			})

			It("returns a GetRequest containing the contents of the request", func() {
				Expect(request).To(BeEquivalentTo(&http.GetRequest{
					BaseDirectory: "/public",
					Target:  "/foo",
					Version: "HTTP/1.1",
				}))
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("reads the entire request until reaching a line with only CRLF", func() {
				Expect(reader.Buffered()).To(Equal(0))
			})
		})

		Context("given any other request method", func() {
			It("returns a ParseError with 501 Not Implemented", func() {
				request, err = parser.ParseRequest(makeReader("get / HTTP/1.1\r\n\n"))
				Expect(err).To(BeEquivalentTo(&http.ParseError{
					StatusCode: 501,
					Reason:     "Not Implemented"}))
			})
		})
	})
})

func makeReader(template string, values ...string) *bufio.Reader {
	text := fmt.Sprintf(template, values)
	return bufio.NewReader(bytes.NewBufferString(text))
}

func beABadRequestError() types.GomegaMatcher {
	return BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"})
}
