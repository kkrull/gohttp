package playground_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRoute", func() {
	It("returns a Route configured for this package", func() {
		route := playground.NewRoute()
		Expect(route).NotTo(BeNil())
		Expect(route.Controller).To(BeAssignableToTypeOf(&playground.StatelessOptionController{}))
	})
})

var _ = Describe("Route", func() {
	Describe("#Route", func() {
		var (
			router        http.Route
			controller    *OptionControllerMock
			requested     *http.RequestLine
			routedRequest http.Request
		)

		BeforeEach(func() {
			controller = &OptionControllerMock{}
			router = &playground.Route{Controller: controller}
		})

		It("routes OPTIONS /method_options", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "/method_options"}
			routedRequest = router.Route(requested)
			routedRequest.Handle(&bufio.Writer{})
			controller.OptionsShouldHaveBeenReceived("/method_options")
		})

		It("routes OPTIONS /method_options2", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "/method_options2"}
			routedRequest = router.Route(requested)
			routedRequest.Handle(&bufio.Writer{})
			controller.OptionsShouldHaveBeenReceived("/method_options2")
		})

		It("returns nil on any other request", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
