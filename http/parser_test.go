package http_test

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/mock"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("RequestLineRouter", func() {
	Describe("#ParseRequest", func() {
		var (
			parser        *http.RequestLineRouter
			request       http.Request
			err           http.Response
			matchAllRoute http.Route
		)

		BeforeEach(func() {
			matchAllRoute = &mock.Route{RouteReturns: mock.Request{}}
		})

		Describe("it returns 400 Bad Request", func() {
			BeforeEach(func() {
				parser = &http.RequestLineRouter{
					Routes: []http.Route{matchAllRoute}}
			})

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

		Context("given a well-formed request", func() {
			var (
				reader *bufio.Reader
			)

			BeforeEach(func() {
				buffer := bytes.NewBufferString("GET /foo HTTP/1.1\r\nAccept: */*\r\n\r\n")
				reader = bufio.NewReader(buffer)
				parser = &http.RequestLineRouter{
					Routes: []http.Route{matchAllRoute}}
				request, err = parser.ParseRequest(reader)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("reads the entire request until reaching a line with only CRLF", func() {
				Expect(reader.Buffered()).To(Equal(0))
			})
		})

		Context("given a well-formed request not matched by any Route", func() {
			It("returns a NotImplemented response", func() {
				parser = &http.RequestLineRouter{}
				request, err = parser.ParseRequest(makeReader("get / HTTP/1.1\r\n\n"))
				Expect(err).To(BeEquivalentTo(&servererror.NotImplemented{Method: "get"}))
			})
		})

		Context("given a well-formed request matched by 1 or more Routes", func() {
			var (
				unrelatedRoute = &mock.Route{RouteReturns: nil}
				matchingRoute  = &mock.Route{RouteReturns: mock.Request{}}
			)

			BeforeEach(func() {
				parser = &http.RequestLineRouter{
					Routes: []http.Route{unrelatedRoute, matchingRoute}}
				request, err = parser.ParseRequest(makeReader("HEAD /foo HTTP/1.1\r\nAccept: */*\r\n\r\n"))
			})

			It("tries routing the method and target from the request until it finds a match", func() {
				unrelatedRoute.ShouldHaveReceived("HEAD", "/foo")
				matchingRoute.ShouldHaveReceived("HEAD", "/foo")
			})

			It("returns the request from the first matching Route", func() {
				Expect(request).To(Equal(matchingRoute.RouteReturns))
				Expect(err).NotTo(HaveOccurred())
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
