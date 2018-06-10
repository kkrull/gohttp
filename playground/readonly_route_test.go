package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewReadOnlyRoute", func() {
	It("returns a route to a resource that only supports read methods", func() {
		route := playground.NewReadOnlyRoute("/readonly")
		Expect(route).NotTo(BeNil())
		Expect(route.Path).To(Equal("/readonly"))
		Expect(route.Resource).To(BeAssignableToTypeOf(&playground.ReadableNopResource{}))
	})
})

var _ = Describe("ReadOnlyRoute", func() {
	const (
		configuredPath  = "/method_options2"
		nonMatchingPath = "/"
	)

	Describe("#Route", func() {
		var (
			router   http.Route
			resource *ReadOnlyResourceMock
		)

		BeforeEach(func() {
			resource = &ReadOnlyResourceMock{}
			router = &playground.ReadOnlyRoute{
				Path:     configuredPath,
				Resource: resource,
			}
		})

		Context("when the path is /method_options2", func() {
			It("routes GET to Teapot#Get", func() {
				handleRequest(router, http.GET, configuredPath)
				resource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Teapot#Options", func() {
				handleRequest(router, http.HEAD, configuredPath)
				resource.HeadShouldHaveBeenCalled()
			})

			Context("when the method is OPTIONS", func() {
				var response = &bytes.Buffer{}

				BeforeEach(func() {
					requested := http.NewOptionsMessage(configuredPath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())

					response.Reset()
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, http.GET, http.HEAD, http.OPTIONS))
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := http.NewPutMessage(configuredPath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed(http.GET, http.HEAD, http.OPTIONS)))
			})
		})

		It("returns nil on any other path", func() {
			requested := http.NewGetMessage(nonMatchingPath)
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})

var _ = Describe("ReadableNopResource", func() {
	var (
		controller *playground.ReadableNopResource
		response   = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
		controller = &playground.ReadableNopResource{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(response, http.NewGetMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response, http.NewHeadMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})
})
