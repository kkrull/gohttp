package teapot_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("teapotRoute", func() {
	var (
		router        http.Route
		controller    *ControllerMock
		requested     *http.RequestLine
		routedRequest http.Request
	)

	BeforeEach(func() {
		controller = &ControllerMock{}
		router = &teapot.Route{Controller: controller}
	})

	Describe("#Route", func() {
		Context("given GET /coffee", func() {
			It("returns a GetRequest", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/coffee"}
				routedRequest = router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(&teapot.GetRequest{}))
			})

			It("routes to Controller#Get", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/coffee"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.GetShouldHaveReceivedTarget("/coffee")
			})
		})

		Context("given GET /tea", func() {
			It("routes to Controller#Get", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/tea"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.GetShouldHaveReceivedTarget("/tea")
			})
		})

		It("passes on any other method", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "/tea"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})

		It("passes on any other target", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
