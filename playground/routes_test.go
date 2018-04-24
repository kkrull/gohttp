package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/playground"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewReadOnlyRoute", func() {
	It("returns a route to a resource that only supports read methods", func() {
		route := playground.NewReadOnlyRoute()
		Expect(route).NotTo(BeNil())
		Expect(route.Resource).To(BeAssignableToTypeOf(&playground.ReadableNopResource{}))
	})
})

var _ = Describe("ReadOnlyRoute", func() {
	Describe("#Route", func() {
		var (
			router           http.Route
			readOnlyResource *ReadOnlyResourceMock
		)

		BeforeEach(func() {
			readOnlyResource = &ReadOnlyResourceMock{}
			router = &playground.ReadOnlyRoute{Resource: readOnlyResource}
		})

		Context("when the target is /method_options2", func() {
			It("routes GET to Resource#Get", func() {
				handleRequest(router, "GET", "/method_options2")
				readOnlyResource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Resource#Options", func() {
				handleRequest(router, "HEAD", "/method_options2")
				readOnlyResource.HeadShouldHaveBeenCalled()
			})

			Context("when the method is OPTIONS", func() {
				var response = &bytes.Buffer{}

				BeforeEach(func() {
					requested := &http.RequestLine{Method: "OPTIONS", Target: "/method_options2"}
					routedRequest := router.Route(requested)
					ExpectWithOffset(1, routedRequest).NotTo(BeNil())

					response.Reset()
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, "GET", "HEAD", "OPTIONS"))
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
			It("routes GET to Resource#Get", func() {
				handleRequest(router, "GET", "/method_options")
				resource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Resource#Head", func() {
				handleRequest(router, "HEAD", "/method_options")
				resource.HeadShouldHaveBeenCalled()
			})

			Context("when the method is OPTIONS", func() {
				var response = &bytes.Buffer{}

				BeforeEach(func() {
					requested := &http.RequestLine{Method: "OPTIONS", Target: "/method_options"}
					routedRequest := router.Route(requested)
					ExpectWithOffset(1, routedRequest).NotTo(BeNil())

					response.Reset()
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, "GET", "HEAD", "OPTIONS", "POST", "PUT"))
			})

			It("routes POST to Resource#Post", func() {
				handleRequest(router, "POST", "/method_options")
				resource.PostShouldHaveBeenCalled()
			})

			It("routes PUT to Resource#Put", func() {
				handleRequest(router, "PUT", "/method_options")
				resource.PutShouldHaveBeenCalled()
			})

			It("returns nil on any other method", func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/method_options"}
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

	routedRequest.Handle(&bytes.Buffer{})
}
