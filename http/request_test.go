package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"bytes"
	"bufio"
	"github.com/kkrull/gohttp/http"
)

var _ = Describe("RFC7230RequestParser", func() {
	Describe("#ParseRequest", func() {
		var (
			parser  *http.RFC7230RequestParser
			reader  *bufio.Reader
			request *http.Request
			err     *http.ParseError
		)

		BeforeEach(func() {
			parser = &http.RFC7230RequestParser{}
		})

		Describe("it returns 400 Bad Request", func() {
			It("for a completely blank request", func() {
				buffer := bytes.NewBufferString("")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("for any line missing CR", func() {
				buffer := bytes.NewBufferString("GET / HTTP/1.1\r\n\n")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("for any line missing LF", func() {
				buffer := bytes.NewBufferString("GET / HTTP/1.1\r")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("for a request missing a request-line", func() {
				buffer := bytes.NewBufferString("\r\n")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("for a request missing an ending CRLF", func() {
				buffer := bytes.NewBufferString("GET / HTTP/1.1\r\n")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("when multiple spaces are separating fields in request-line", func() {
				buffer := bytes.NewBufferString("GET /  HTTP/1.1\r\n\r\n")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("when fields in request-line contain spaces", func() {
				buffer := bytes.NewBufferString("GET /a\\ b HTTP/1.1\r\n\r\n")
				request, err = parser.ParseRequest(bufio.NewReader(buffer))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})

			It("when the request starts with whitespace", func() {
				request, err = parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
				Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
			})
		})

		Context("given a well-formed request", func() {
			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)
				request, err = parser.ParseRequest(reader)
			})

			It("parses the request", func() {
				Expect(request).To(BeEquivalentTo(&http.Request{
					Method:  "GET",
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
	})
})

func makeReader(template string, values ...string) *bufio.Reader {
	text := fmt.Sprintf(template, values)
	return bufio.NewReader(bytes.NewBufferString(text))
}
