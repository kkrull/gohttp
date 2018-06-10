package capability_test

import (
	"bufio"

	"github.com/kkrull/gohttp/capability"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const serverCapabilityTarget = "*"

var _ = Describe("::NewRoute", func() {
	It("configures the route with StaticCapabilityServer", func() {
		route := capability.NewRoute(serverCapabilityTarget)
		Expect(route.Controller).To(BeAssignableToTypeOf(&capability.StaticCapabilityServer{}))
		Expect(route.Target).To(Equal(serverCapabilityTarget))
	})

	It("configures available methods to the server as GET and HEAD", func() {
		route := capability.NewRoute(serverCapabilityTarget)
		Expect(route.Controller).To(BeEquivalentTo(
			&capability.StaticCapabilityServer{
				AvailableMethods: []string{http.GET, http.HEAD},
			},
		))
	})
})

var _ = Describe("ServerCapabilityRoute", func() {
	Describe("#Route", func() {
		var (
			router        http.Route
			controller    *ServerCapabilityServerMock
			requested     http.RequestMessage
			routedRequest http.Request
		)

		BeforeEach(func() {
			controller = &ServerCapabilityServerMock{}
			router = &capability.ServerCapabilityRoute{
				Target:     serverCapabilityTarget,
				Controller: controller,
			}
		})

		Context("when the path is the server capability target (*)", func() {
			It("routes OPTIONS to ServerResource", func() {
				requested = http.NewOptionsMessage("*")
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.OptionsShouldHaveBeenCalled()
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested = http.NewGetMessage(serverCapabilityTarget)
				routedRequest = router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed(http.OPTIONS)))
			})
		})

		It("returns nil to pass on any other path", func() {
			requested = http.NewOptionsMessage("/something/else")
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
