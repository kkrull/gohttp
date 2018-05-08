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
			Reporter: &playground.AssignmentReporter{},
		}))
	})
})

var _ = Describe("ParameterRoute", func() {
	Describe("#Route", func() {
		var (
			router   http.Route
			reporter *ParameterReporterMock
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			reporter = &ParameterReporterMock{}
			router = &playground.ParameterRoute{Reporter: reporter}
			response.Reset()
		})

		Context("when the path is /parameters", func() {
			It("routes GET to ParameterReporter#Get with the decoded query parameters", func() {
				request := &httptest.RequestMock{}
				requested := &httptest.RequestMessage{
					PathReturns:                "/parameters",
					MakeResourceRequestReturns: request,
				}

				Expect(router.Route(requested)).To(BeIdenticalTo(request))
				requested.MakeResourceRequestShouldHaveReceived(reporter)
			})

			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage("/parameters")
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("sets Allow to the methods implemented by this type",
					httptest.ShouldAllowMethods(response, http.GET, http.OPTIONS))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage("/parameters")
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

var _ = Describe("AssignmentReporter", func() {
	Describe("#Get", func() {
		var (
			controller      *playground.AssignmentReporter
			request         *httptest.RequestMessage
			responseMessage *httptest.ResponseMessage

			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			response.Reset()
		})

		Context("given any number of query parameters", func() {
			BeforeEach(func() {
				controller = &playground.AssignmentReporter{}
				request = &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   "/parameters",
				}

				controller.Get(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("responds 200 OK", func() {
				responseMessage.StatusShouldBe(200, "OK")
				responseMessage.ShouldBeWellFormed()
			})
			It("sets Content-Type to text/plain", func() {
				responseMessage.HeaderShould("Content-Type", Equal("text/plain"))
			})
		})

		Context("given a request with no query parameters", func() {
			BeforeEach(func() {
				controller = &playground.AssignmentReporter{}
				request = &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   "/parameters",
				}

				controller.Get(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("sets Content-Length to 0", func() {
				responseMessage.HeaderShould("Content-Length", Equal("0"))
			})
			It("writes no body", func() {
				responseMessage.BodyShould(BeEmpty())
			})
		})

		Context("given a request with 1 or more query parameters", func() {
			BeforeEach(func() {
				controller = &playground.AssignmentReporter{}
				request = &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   "/parameters",
				}
				request.AddQueryParameter("foo", "bar")

				controller.Get(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("sets Content-Length", func() {
				responseMessage.HeaderShould("Content-Length", Not(Equal("0")))
			})
			It("lists each parameter and its value in the body", func() {
				responseMessage.BodyShould(ContainSubstring("foo = bar"))
			})
		})
	})
})
