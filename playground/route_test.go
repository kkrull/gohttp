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
		Expect(route.Controller).To(BeAssignableToTypeOf(&playground.AllowedMethodsController{}))
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
			router = &playground.Route{Controller: controller}
		})

		It("routes GET /method_options to Controller#Get", func() {
			handleRequest(router, "GET", "/method_options")
			controller.GetShouldHaveBeenReceived("/method_options")
		})

		It("routes HEAD /method_options to Controller#Head", func() {
			handleRequest(router, "HEAD", "/method_options")
			controller.HeadShouldHaveBeenReceived("/method_options")
		})

		It("routes OPTIONS /method_options to Controller#Options", func() {
			handleRequest(router, "OPTIONS", "/method_options")
			controller.OptionsShouldHaveBeenReceived("/method_options")
		})

		It("routes POST /method_options to Controller#Post", func() {
			handleRequest(router, "POST", "/method_options")
			controller.PostShouldHaveBeenReceived("/method_options")
		})

		It("routes PUT /method_options to Controller#Put", func() {
			handleRequest(router, "PUT", "/method_options")
			controller.PutShouldHaveBeenReceived("/method_options")
		})

		It("routes OPTIONS /method_options2 to Controller#Options", func() {
			handleRequest(router, "OPTIONS", "/method_options2")
			controller.OptionsShouldHaveBeenReceived("/method_options2")
		})

		It("returns nil on any other request", func() {
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
