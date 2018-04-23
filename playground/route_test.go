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
		Expect(route.WriteController).To(BeAssignableToTypeOf(&playground.WritableNopController{}))
	})
})

var _ = Describe("Route", func() {
	Describe("#Route", func() {
		var (
			router     http.Route
			controller *ControllerMock
		)

		BeforeEach(func() {
			controller = &ControllerMock{}
			router = &playground.Route{WriteController: controller}
		})

		Context("when the target is /method_options", func() {
			It("routes GET to WriteController#Get", func() {
				handleRequest(router, "GET", "/method_options")
				controller.GetShouldHaveBeenReceived("/method_options")
			})

			It("routes HEAD to WriteController#Head", func() {
				handleRequest(router, "HEAD", "/method_options")
				controller.HeadShouldHaveBeenReceived("/method_options")
			})

			It("routes OPTIONS to WriteController#Options", func() {
				handleRequest(router, "OPTIONS", "/method_options")
				controller.OptionsShouldHaveBeenReceived("/method_options")
			})

			It("routes POST to WriteController#Post", func() {
				handleRequest(router, "POST", "/method_options")
				controller.PostShouldHaveBeenReceived("/method_options")
			})

			It("routes PUT to WriteController#Put", func() {
				handleRequest(router, "PUT", "/method_options")
				controller.PutShouldHaveBeenReceived("/method_options")
			})
		})

		Context("when the target is /method_options2", func() {
			It("routes GET to WriteController#Get", func() {
				handleRequest(router, "GET", "/method_options2")
				controller.GetShouldHaveBeenReceived("/method_options2")
			})

			It("routes HEAD to WriteController#Options", func() {
				handleRequest(router, "HEAD", "/method_options2")
				controller.HeadShouldHaveBeenReceived("/method_options2")
			})

			It("routes OPTIONS to WriteController#Options", func() {
				handleRequest(router, "OPTIONS", "/method_options2")
				controller.OptionsShouldHaveBeenReceived("/method_options2")
			})

			It("returns nil for any other method", func() {
				requested := &http.RequestLine{Method: "PUT", Target: "/method_options2"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeNil())
			})
		})

		It("returns nil on any other target", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/"}
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})

func handleRequest(router http.Route, method, target string) {
	requested := &http.RequestLine{Method: method, Target: target}
	routedRequest := router.Route(requested)
	ExpectWithOffset(1, routedRequest).NotTo(BeNil())
	routedRequest.Handle(&bufio.Writer{})
}
