package capability_test

import (
	"bufio"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRoute", func() {
	It("configures the route with StaticCapabilityServer", func() {
		route := capability.NewRoute()
		Expect(route.Controller).To(BeAssignableToTypeOf(&capability.StaticCapabilityServer{}))
	})

	It("configures available methods to the server as GET and HEAD", func() {
		route := capability.NewRoute()
		Expect(route.Controller).To(BeEquivalentTo(
			&capability.StaticCapabilityServer{
				AvailableMethods: []string{"GET", "HEAD"},
			},
		))
	})
})

var _ = Describe("ServerCapabilityRoute", func() {
	Describe("#Route", func() {
		var (
			router        http.Route
			controller    *ServerCapabilityServerMock
			requested     *http.RequestLine
			routedRequest http.Request
		)

		BeforeEach(func() {
			controller = &ServerCapabilityServerMock{}
			router = &capability.ServerCapabilityRoute{Controller: controller}
		})

		Context("when the target is *", func() {
			It("routes OPTIONS to ServerResource", func() {
				requested = &http.RequestLine{Method: "OPTIONS", Target: "*"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.OptionsShouldHaveBeenCalled()
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested = &http.RequestLine{Method: "GET", Target: "*"}
				routedRequest = router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("OPTIONS")))
			})
		})

		It("returns nil to pass on any other target", func() {
			requested = &http.RequestLine{Method: "OPTIONS", Target: "/"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
