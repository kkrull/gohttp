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
	Describe("#RouteRequest", func() {
		var (
			router        *http.RequestLineRouter
			request       http.Request
			err           http.Response
			matchAllRoute *mock.Route
		)

		BeforeEach(func() {
			matchAllRoute = &mock.Route{RouteReturns: &mock.Request{}}
		})

		Context("given a well-formed request not matched by any Route", func() {
			It("returns a NotImplemented response", func() {
				router = http.NewRouter()
				request, err = router.RouteRequest(makeReader("get / HTTP/1.1\r\n\r\n"))
				Expect(err).To(BeEquivalentTo(&servererror.NotImplemented{Method: "get"}))
			})
		})

		Context("given a well-formed request matched by 1 or more Routes", func() {
			var (
				unrelatedRoute = &mock.Route{RouteReturns: nil}
				matchingRoute  = &mock.Route{RouteReturns: &mock.Request{}}
			)

			BeforeEach(func() {
				router = http.NewRouter()
				router.AddRoute(unrelatedRoute)
				router.AddRoute(matchingRoute)
				request, err = router.RouteRequest(makeReader("HEAD /foo HTTP/1.1\r\nAccept: */*\r\n\r\n"))
			})

			It("tries routing the method and path from the request until it finds a match", func() {
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

func makeReader(template string, values ...interface{}) *bufio.Reader {
	text := fmt.Sprintf(template, values...)
	return bufio.NewReader(bytes.NewBufferString(text))
}

func beABadRequestResponse(why string) types.GomegaMatcher {
	return BeEquivalentTo(&clienterror.BadRequest{DisplayText: why})
}
