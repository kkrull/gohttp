package capability_test

import (
	"bufio"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//Panics due to incomplete wiring are easy to cause due to Go's permissive struct declaration
//and hard to root cause when the server goes down and starts refusing connections behind a FitNesse suite
//that swallows console output from the server
var _ = Describe("::NewRoute", func() {
	It("configures the route with StaticCapabilityController", func() {
		route := capability.NewRoute()
		Expect(route.Controller).To(BeAssignableToTypeOf(&capability.StaticCapabilityController{}))
	})

	It("configures available methods to the server as GET and HEAD", func() {
		route := capability.NewRoute()
		Expect(route.Controller).To(BeEquivalentTo(
			&capability.StaticCapabilityController{
				AvailableMethods: []string{"GET", "HEAD"},
			},
		))
	})
})

var _ = Describe("ServerCapabilityRoute", func() {
	Describe("#Route", func() {
		var (
			router        http.Route
			controller    *ServerCapabilityControllerMock
			requested     *http.RequestLine
			routedRequest http.Request
		)

		BeforeEach(func() {
			controller = &ServerCapabilityControllerMock{}
			router = &capability.ServerCapabilityRoute{Controller: controller}
		})

		It("routes OPTIONS * to ServerCapabilityController#Options", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "*"}
			routedRequest = router.Route(requested)
			routedRequest.Handle(&bufio.Writer{})
			controller.OptionsShouldHaveBeenCalled()
		})

		It("returns nil to pass on any other method", func() {
			requested = &http.RequestLine{Method: "GET", Target: "*"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})

		It("returns nil to pass on any other target", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "/"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
