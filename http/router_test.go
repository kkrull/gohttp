package http_test

import (
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/msg/servererror"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRouter", func() {
	It("configures no logging for requests, by default", func() {
		router := http.NewRouter()
		router.RouteRequest(makeReader("GET / HTTP/1.1\r\n\r\n"))
	})
})

var _ = Describe("RequestLineRouter", func() {
	Describe("#RouteRequest", func() {
		var (
			router  *http.RequestLineRouter
			logger  *RequestLoggerMock
			request http.Request
			err     http.Response
		)

		It("logs the parsed request", func() {
			router = http.NewRouter()
			logger = &RequestLoggerMock{}
			router.LogRequests(logger)
			request, err = router.RouteRequest(makeReader("GET / HTTP/1.1\r\n\r\n"))
			logger.ParsedShouldHaveReceived(http.GET, "/")
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
				unrelatedRoute = &RouteMock{RouteReturns: nil}
				matchingRoute  = &RouteMock{RouteReturns: &httptest.RequestMock{}}
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
