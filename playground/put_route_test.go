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

var _ = Describe("::NewPuttableRoute", func() {
	It("returns a PuttableRoute at the given path", func() {
		route := playground.NewPuttableRoute("/oracle")
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&playground.PuttableRoute{
			Path:     "/oracle",
			Resource: &playground.PuttableResource{},
		}))
	})
})

var _ = Describe("PuttableRoute", func() {
	Describe("#Route", func() {
		const givenPath = "/sweetness"

		var (
			router   http.Route
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			router = &playground.PuttableRoute{Path: givenPath}
			response.Reset()
		})

		Context("when the path is exactly equal to the given path", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage(givenPath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows OPTIONS", httptest.ShouldAllowMethods(response, http.OPTIONS))
				It("allows POST", httptest.ShouldAllowMethods(response, http.POST))
				It("allows PUT", httptest.ShouldAllowMethods(response, http.PUT))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage(givenPath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		It("returns nil for any other path", func() {
			requested := http.NewGetMessage("/")
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})

var _ = Describe("PuttableResource", func() {
	var (
		okResource      *playground.PuttableResource
		request         *httptest.RequestMessage
		responseMessage *httptest.ResponseMessage

		response = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
	})

	Describe("#Post", func() {
		Context("given any data in the body", func() {
			BeforeEach(func() {
				okResource = &playground.PuttableResource{}
				request = &httptest.RequestMessage{
					MethodReturns: http.POST,
					PathReturns:   "/form",
				}

				okResource.Post(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("responds 200 OK", func() {
				responseMessage.ShouldBeWellFormed()
				responseMessage.StatusShouldBe(200, "OK")
			})
		})
	})

	Describe("#Put", func() {
		Context("given any data in the body", func() {
			BeforeEach(func() {
				okResource = &playground.PuttableResource{}
				request = &httptest.RequestMessage{
					MethodReturns: http.PUT,
					PathReturns:   "/form",
				}

				okResource.Put(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("responds 200 OK", func() {
				responseMessage.ShouldBeWellFormed()
				responseMessage.StatusShouldBe(200, "OK")
			})
		})
	})
})
