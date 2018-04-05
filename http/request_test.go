package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"bytes"
	"bufio"
	"github.com/kkrull/gohttp/http"
	"strings"
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

		It("returns 400 Bad Request for a completely blank request", func() {
			buffer := bytes.NewBufferString("")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request for any line missing CR", func() {
			buffer := bytes.NewBufferString("GET / HTTP/1.1\r\n\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request for any line missing LF", func() {
			buffer := bytes.NewBufferString("GET / HTTP/1.1\r")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request for a request missing a request-line", func() {
			buffer := bytes.NewBufferString("\r\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request for a request missing an ending CRLF", func() {
			buffer := bytes.NewBufferString("GET / HTTP/1.1\r\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request when multiple spaces are separating fields in request-line", func() {
			buffer := bytes.NewBufferString("GET /  HTTP/1.1\r\n\r\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request when fields in request-line contain spaces", func() {
			buffer := bytes.NewBufferString("GET /a\\ b HTTP/1.1\r\n\r\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
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

	XDescribe("#OldParseRequest", func() {
		var (
			parser  http.RequestParser
			request *http.Request
			err     *http.ParseError
		)

		BeforeEach(func() {
			parser = http.RFC7230RequestParser{}
		})

		Describe("Section 3.1.1", func() {
			Context("given a well formed request line that is not too long", func() {
				It("parses the HTTP method as all text before the first space", func() {
					request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n\r\n"))
					Expect(request).To(BeEquivalentTo(&http.Request{
						Method:  "GET",
						Target:  "/",
						Version: "HTTP/1.1",
					}))
				})
			})

			Context("when the request starts with whitespace", func() {
				It("Returns 400 Bad Request", func() {
					request, err = parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
				})
			})

			Context("when the request method is longer than 8,000 octets", func() {
				It("Returns 501 Not Implemented", func() {
					enormousMethod := strings.Repeat("POST", 2000) + "!"
					request, err = parser.ParseRequest(makeReader("%s / HTTP/1.1\r\n\r\n", enormousMethod))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 501, Reason: "Not Implemented"}))
				})
			})

			Context("when the request target is longer than 8,000 octets", func() {
				It("Returns 414 URI Too Long", func() {
					enormousTarget := strings.Repeat("/foo", 2000) + "/"
					request, err = parser.ParseRequest(makeReader("GET %s HTTP/1.1\r\n\r\n", enormousTarget))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 414, Reason: "URI Too Long"}))
				})
			})
		})
	})
})

func makeReader(template string, values ...string) *bufio.Reader {
	text := fmt.Sprintf(template, values)
	return bufio.NewReader(bytes.NewBufferString(text))
}
