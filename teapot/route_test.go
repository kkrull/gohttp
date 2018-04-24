package teapot_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("teapotRoute", func() {
	var (
		router        http.Route
		controller    *TeapotMock
		requested     *http.RequestLine
		routedRequest http.Request
	)

	BeforeEach(func() {
		controller = &TeapotMock{}
		router = &teapot.Route{Resource: controller}
	})

	Describe("#Route", func() {
		Context("when the target is /coffee", func() {
			It("routes GET /coffee to Resource#GetCoffee", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/coffee"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.GetShouldHaveReceived("/coffee")
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/coffee"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "OPTIONS")))
			})
		})

		Context("when the target is /tea", func() {
			It("routes GET /tea to Resource#GetCoffee", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/tea"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.GetShouldHaveReceived("/tea")
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/tea"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "OPTIONS")))
			})
		})

		It("passes on any other target", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
