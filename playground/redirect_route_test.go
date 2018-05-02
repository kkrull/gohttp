package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/mock"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRedirectRoute", func() {
	It("returns a fully configured RedirectRoute", func() {
		route := playground.NewRedirectRoute()
		Expect(route).To(BeAssignableToTypeOf(&playground.RedirectRoute{}))
		Expect(route).To(BeEquivalentTo(&playground.RedirectRoute{
			Resource: &playground.GoBackHomeResource{},
		}))
	})
})

var _ = Describe("RedirectRoute", func() {
	Describe("#Route", func() {
		var (
			router   http.Route
			resource *playground.GoBackHomeResource
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			resource = &playground.GoBackHomeResource{}
			router = &playground.RedirectRoute{Resource: resource}
			response.Reset()
		})

		Context("when the path is /redirect", func() {
			It("routes GET to RelocatedResource#Get", func() {
				requestMessage := &httptest.RequestMessage{
					PathReturns:                "/redirect",
					MakeResourceRequestReturns: &mock.Request{},
				}

				Expect(router.Route(requestMessage)).To(BeIdenticalTo(requestMessage.MakeResourceRequestReturns))
			})

			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage("/redirect")
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, "GET", "OPTIONS"))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage("/redirect")
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		It("returns nil for any other path", func() {
			requestMessage := http.NewGetMessage("/no-route-for-you")
			Expect(router.Route(requestMessage)).To(BeNil())
		})
	})
})
