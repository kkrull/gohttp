package http_test

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/mock"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("RFC7230RequestParser", func() {
	Describe("#ParseRequest", func() {
		var (
			parser  *http.RFC7230RequestParser
			request http.Request
			err     http.Response
		)

		BeforeEach(func() {
			parser = &http.RFC7230RequestParser{BaseDirectory: "/tmp"}
		})

		Describe("it returns 400 Bad Request", func() {
			It("for a completely blank request", func() {
				request, err = parser.ParseRequest(makeReader(""))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("for any line missing CR", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n\n"))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("for any line missing LF", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r"))
				Expect(err).To(beABadRequestResponse("message header line does not end in LF"))
			})

			It("for a request missing a request-line", func() {
				request, err = parser.ParseRequest(makeReader("\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("for a request missing an ending CRLF", func() {
				request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n"))
				Expect(err).To(beABadRequestResponse("line in request header not ending in CRLF"))
			})

			It("when multiple spaces are separating fields in request-line", func() {
				request, err = parser.ParseRequest(makeReader("GET /  HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("when fields in request-line contain spaces", func() {
				request, err = parser.ParseRequest(makeReader("GET /a\\ b HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})

			It("when the request starts with whitespace", func() {
				request, err = parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
				Expect(err).To(beABadRequestResponse("incorrectly formatted or missing request-line"))
			})
		})

		Context("given a well-formed request that no Route can handle", func() {
			XIt("returns a NotImplemented response")
		})

		Context("given a well-formed request", func() {
			var (
				matchingRoute *mock.Route
			)

			BeforeEach(func() {
				matchingRoute = &mock.Route{}
				parser = &http.RFC7230RequestParser{
					BaseDirectory: "/tmp",
					Routes:        []http.Route{matchingRoute}}
			})

			It("delegates the request to a Route", func() {
				request, err = parser.ParseRequest(makeReader("HEAD /foo HTTP/1.1\r\nAccept: */*\r\n\r\n"))
				matchingRoute.ShouldHaveReceived("HEAD", "/foo")
			})
		})

		Context("given a well-formed GET request", func() {
			var (
				reader       *bufio.Reader
				typedRequest *fs.GetRequest
			)

			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)

				parser = &http.RFC7230RequestParser{BaseDirectory: "/public"}
				request, err = parser.ParseRequest(reader)
				typedRequest = request.(*fs.GetRequest)
			})

			It("returns a GetRequest containing the contents of the request", func() {
				Expect(request).To(BeEquivalentTo(&fs.GetRequest{
					BaseDirectory: "/public",
					Target:        "/foo",
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
			It("returns a NotImplemented response", func() {
				request, err = parser.ParseRequest(makeReader("get / HTTP/1.1\r\n\n"))
				Expect(err).To(BeEquivalentTo(&servererror.NotImplemented{Method: "get"}))
			})
		})
	})
})

func makeReader(template string, values ...string) *bufio.Reader {
	text := fmt.Sprintf(template, values)
	return bufio.NewReader(bytes.NewBufferString(text))
}

func beABadRequestResponse(why string) types.GomegaMatcher {
	return BeEquivalentTo(&clienterror.BadRequest{DisplayText: why})
}
