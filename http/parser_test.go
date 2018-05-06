package http_test

import (
	"bufio"
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/servererror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LineRequestParser", func() {
	Describe("#Parse", func() {
		var (
			parser  http.RequestParser
			request http.RequestMessage
			err     http.Response
		)

		BeforeEach(func() {
			parser = &http.LineRequestParser{}
		})

		Describe("it returns 400 Bad Request", func() {
			It("for a completely blank request", func() {
				request, err = parser.Parse(makeReader(""))
				Expect(err).To(beABadRequestResponse("end of input before terminating CRLF"))
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
				Expect(err).To(beABadRequestResponse("end of input before terminating CRLF"))
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

		Context("given a request with an unknown method", func() {
			It("returns a NotImplemented error response", func() {
				//RFC 7231 Section 4.1
				requestWithFieldsOutOfOrder := "/ GET HTTP/1.1\r\n\r\n"
				request, err = parser.Parse(makeReader(requestWithFieldsOutOfOrder))
				Expect(err).To(BeEquivalentTo(&servererror.NotImplemented{Method: "/"}))
				Expect(err).NotTo(BeNil())
			})
		})

		Context("given a well-formed request with a known method", func() {
			var reader *bufio.Reader

			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)
				request, err = parser.Parse(reader)
			})

			It("parses the request-line", func() {
				Expect(request.Method()).To(Equal(http.GET))
				Expect(request.Target()).To(Equal("/foo"))
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("reads the entire request until reaching a line with only CRLF", func() {
				Expect(reader.Buffered()).To(Equal(0))
			})
		})

		Context("given a target with no query or fragment", func() {
			BeforeEach(func() {
				request, _ = parser.Parse(requestWithTarget("/widget"))
			})

			It("the path is the full target", func() {
				Expect(request.Path()).To(Equal("/widget"))
			})
			It("there are no query parameters", func() {
				Expect(request.QueryParameters()).To(BeEmpty())
			})
		})

		Context("given a target with a query", func() {
			It("the path is the part before the ?", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?field=value"))
				Expect(request.Path()).To(Equal("/widget"))
			})

			It("parses the part after the ? into query parameters", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?field=value"))
				Expect(request.QueryParameters()).To(
					ContainElement(http.QueryParameter{Name: "field", Value: "value"}))
			})

			It("parses parameters without a value into QueryParameter#Name", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?flag"))
				Expect(request.QueryParameters()).To(
					ContainElement(http.QueryParameter{Name: "flag", Value: ""}))
			})

			It("uses '=' to split a parameter's name and value", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?field=value"))
				Expect(request.QueryParameters()).To(
					ContainElement(http.QueryParameter{Name: "field", Value: "value"}))
			})

			It("uses '&' to split among multiple parameters", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?one=1&two=2"))
				Expect(request.QueryParameters()).To(Equal([]http.QueryParameter{
					{Name: "one", Value: "1"},
					{Name: "two", Value: "2"},
				}))
			})
		})

		Context("given a target with a fragment", func() {
			BeforeEach(func() {
				request, _ = parser.Parse(requestWithTarget("/widget#section"))
			})

			It("the path is the part before the '#'", func() {
				Expect(request.Path()).To(Equal("/widget"))
			})
			It("there are no query parameters", func() {
				Expect(request.QueryParameters()).To(BeEmpty())
			})
		})

		Context("given a target with a query and a fragment", func() {
			BeforeEach(func() {
				request, _ = parser.Parse(requestWithTarget("/widget?field=value#section"))
			})

			It("the path is the part before the ?", func() {
				Expect(request.Path()).To(Equal("/widget"))
			})
			It("query parameters are parsed from the part between the ? and the #", func() {
				Expect(request.QueryParameters()).To(Equal([]http.QueryParameter{
					{Name: "field", Value: "value"},
				}))
			})
		})

		Context("given a target with percent-encoded query parameters", func() {
			It("decodes percent-encoded values", func() {
				request, _ = parser.Parse(requestWithTarget("/widget?less=%3C"))
				Expect(request.QueryParameters()).To(Equal([]http.QueryParameter{
					{Name: "less", Value: "<"},
				}))
			})
		})
	})
})

func requestWithTarget(target string) *bufio.Reader {
	return makeReader("GET %s HTTP/1.1\r\n\r\n", target)
}
