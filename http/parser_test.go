package http_test

import (
	"bufio"
	"bytes"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LineRequestParser", func() {
	Describe("#Parse", func() {
		var (
			parser  http.RequestParser
			request *http.RequestLine
			err     http.Response
		)

		BeforeEach(func() {
			parser = &http.LineRequestParser{}
		})

		Describe("it returns 400 Bad Request", func() {
			It("for a completely blank request", func() {
				request, err = parser.Parse(makeReader(""))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("for any line missing CR", func() {
				request, err = parser.Parse(makeReader("GET / HTTP/1.1\r\n\n"))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("for any line missing LF", func() {
				request, err = parser.Parse(makeReader("GET / HTTP/1.1\r"))
				Expect(err).To(beABadRequestResponse("message header line does not end in LF"))
			})

			It("for a request missing a request-line", func() {
				request, err = parser.Parse(makeReader("\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("for a request missing an ending CRLF", func() {
				request, err = parser.Parse(makeReader("GET / HTTP/1.1\r\n"))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("when multiple spaces are separating fields in request-line", func() {
				request, err = parser.Parse(makeReader("GET /  HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("when fields in request-line contain spaces", func() {
				request, err = parser.Parse(makeReader("GET /a\\ b HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("when the request starts with whitespace", func() {
				request, err = parser.Parse(makeReader(" GET / HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})
		})

		Context("given a well-formed request", func() {
			var (
				reader *bufio.Reader
			)

			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)
				request, err = parser.Parse(reader)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("reads the entire request until reaching a line with only CRLF", func() {
				Expect(reader.Buffered()).To(Equal(0))
			})
		})

		Context("given a well-formed request with query parameters", func() {
			BeforeEach(func() {
				request, err = parser.Parse(makeReader("GET /foo?one=1 HTTP/1.1\r\n\r\n"))
				Expect(err).NotTo(HaveOccurred())
			})

			XIt("removes the query string from the target")
			XIt("passes decoded parameters to the routes")
		})
	})
})
