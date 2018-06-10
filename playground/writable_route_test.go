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
		route := playground.NewReadWriteRoute("/rw")
		Expect(route).NotTo(BeNil())
		Expect(route.Path).To(Equal("/rw"))
		Expect(route.Resource).To(BeAssignableToTypeOf(&playground.ReadWriteNopResource{}))
	})
})

var _ = Describe("ReadWriteRoute", func() {
	const (
		configuredPath  = "/method_options"
		nonMatchingPath = "/"
	)

	Describe("#Route", func() {
		var (
			router   http.Route
			resource *ReadWriteResourceMock
		)

		BeforeEach(func() {
			resource = &ReadWriteResourceMock{}
			router = &playground.ReadWriteRoute{
				Path:     configuredPath,
				Resource: resource,
			}
		})

		Context("when the path is /method_options", func() {
			It("routes GET to Teapot#Get", func() {
				handleRequest(router, http.GET, configuredPath)
				resource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Teapot#Head", func() {
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
					httptest.ShouldAllowMethods(response, http.GET, http.HEAD, http.OPTIONS, http.POST, http.PUT))
			})

			It("routes POST to Teapot#Post", func() {
				handleRequest(router, http.POST, configuredPath)
				resource.PostShouldHaveBeenCalled()
			})

			It("routes PUT to Teapot#Put", func() {
				handleRequest(router, http.PUT, configuredPath)
				resource.PutShouldHaveBeenCalled()
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := http.NewTraceMessage(configuredPath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed(http.GET, http.HEAD, http.OPTIONS, http.POST, http.PUT)))
			})
		})

		It("returns nil on any other path", func() {
			requested := http.NewGetMessage(nonMatchingPath)
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
			controller.Head(response, http.NewHeadMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(response, http.NewPostMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(response, http.NewPutMessage("/"))
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})
})
