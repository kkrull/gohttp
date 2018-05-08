package http_test

import (
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
			router  *http.RequestLineRouter
			request http.Request
			err     http.Response
		)

		Context("given a well-formed request not matched by any Route", func() {
			It("returns a NotImplemented response", func() {
				router = http.NewRouter()
				request, err = router.RouteRequest(makeReader("get / HTTP/1.1\r\n\r\n"))
				Expect(err).To(BeEquivalentTo(&servererror.NotImplemented{Method: "get"}))
			})
		})

		Context("given a well-formed request matched by 1 or more Routes", func() {
			var (
				unrelatedRoute = &RouteMock{RouteReturns: nil}
				matchingRoute  = &RouteMock{RouteReturns: &mock.Request{}}
			)

			BeforeEach(func() {
				router = http.NewRouter()
				router.AddRoute(unrelatedRoute)
				router.AddRoute(matchingRoute)
				request, err = router.RouteRequest(makeReader("HEAD /foo HTTP/1.1\r\nAccept: */*\r\n\r\n"))
			})

			It("tries routing the method and path from the request until it finds a match", func() {
				unrelatedRoute.ShouldHaveReceived(http.HEAD, "/foo")
				matchingRoute.ShouldHaveReceived(http.HEAD, "/foo")
			})

			It("returns the request from the first matching Route", func() {
				Expect(request).To(Equal(matchingRoute.RouteReturns))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func beABadRequestResponse(why string) types.GomegaMatcher {
	return BeEquivalentTo(&clienterror.BadRequest{DisplayText: why})
}
