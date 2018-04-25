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

var _ = Describe("::NewParameterRoute", func() {
	It("returns a fully configured ParameterRoute", func() {
		route := playground.NewParameterRoute()
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&playground.ParameterRoute{
			Decoder: &playground.TheDecoder{},
		}))
	})
})

var _ = Describe("ParameterRoute", func() {
	Describe("#Route", func() {
		var (
			router   http.Route
			decoder  *ParameterDecoderMock
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			decoder = &ParameterDecoderMock{}
			router = &playground.ParameterRoute{Decoder: decoder}
			response.Reset()
		})

		Context("when the target is /parameters", func() {
			XIt("routes GET to ParameterDecoder#Get with the decoded query parameters", func() {
				requested := &http.RequestLine{
					Method:      "GET",
					Target:      "/parameters",
					QueryString: "one=1&two=2",
				}

				routedRequest := router.Route(requested)
				Expect(routedRequest).NotTo(BeNil())
				routedRequest.Handle(response)

				expectedDecodedParameters := map[string]string{
					"one": "1",
					"two": "2",
				}
				decoder.GetShouldHaveReceived(expectedDecodedParameters)
			})

			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := &http.RequestLine{Method: "OPTIONS", Target: "/parameters"}
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, "GET", "OPTIONS"))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/parameters"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		It("returns nil for any other target", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/"}
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})
