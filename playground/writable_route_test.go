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

var _ = Describe("::NewReadWriteRoute", func() {
	It("returns a route to a resource that supports read and write methods", func() {
		route := playground.NewReadWriteRoute()
		Expect(route).NotTo(BeNil())
		Expect(route.Resource).To(BeAssignableToTypeOf(&playground.ReadWriteNopResource{}))
	})
})

var _ = Describe("ReadWriteRoute", func() {
	Describe("#Route", func() {
		var (
			router   http.Route
			resource *ReadWriteResourceMock
		)

		BeforeEach(func() {
			resource = &ReadWriteResourceMock{}
			router = &playground.ReadWriteRoute{Resource: resource}
		})

		Context("when the target is /method_options", func() {
			It("routes GET to Teapot#Get", func() {
				handleRequest(router, "GET", "/method_options")
				resource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Teapot#Head", func() {
				handleRequest(router, "HEAD", "/method_options")
				resource.HeadShouldHaveBeenCalled()
			})

			Context("when the method is OPTIONS", func() {
				var response = &bytes.Buffer{}

				BeforeEach(func() {
					requested := http.NewOptionsMessage("/method_options")
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())

					response.Reset()
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, "GET", "HEAD", "OPTIONS", "POST", "PUT"))
			})

			It("routes POST to Teapot#Post", func() {
				handleRequest(router, "POST", "/method_options")
				resource.PostShouldHaveBeenCalled()
			})

			It("routes PUT to Teapot#Put", func() {
				handleRequest(router, "PUT", "/method_options")
				resource.PutShouldHaveBeenCalled()
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := http.NewTraceMessage("/method_options")
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "HEAD", "OPTIONS", "POST", "PUT")))
			})
		})

		It("returns nil on any other target", func() {
			requested := http.NewGetMessage("/")
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})

var _ = Describe("ReadWriteNopResource", func() {
	var (
		controller *playground.ReadWriteNopResource
		response   = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
		controller = &playground.ReadWriteNopResource{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(response, http.NewGetMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response, "/")
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(response, "/")
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(response, "/")
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})
})
