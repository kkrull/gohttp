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

var _ = Describe("::NewSingletonRoute", func() {
	It("returns a SingletonRoute at the given path", func() {
		route := playground.NewSingletonRoute("/singleton")
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&playground.SingletonRoute{
			Singleton: &playground.SingletonResource{Path: "/singleton"},
		}))
	})
})

var _ = Describe("SingletonRoute", func() {
	Describe("#Route", func() {
		const givenPath = "/reverie"

		var (
			router   http.Route
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			router = &playground.SingletonRoute{
				Singleton: &playground.SingletonResource{Path: givenPath},
			}

			response.Reset()
		})

		Context("when the path is the configured path", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage(givenPath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows OPTIONS", httptest.ShouldAllowMethods(response, http.OPTIONS))
				It("allows GET", httptest.ShouldAllowMethods(response, http.GET))
				It("allows POST", httptest.ShouldAllowMethods(response, http.POST))
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

var _ = Describe("SingletonResource", func() {
	var (
		singleton       *playground.SingletonResource
		request         *httptest.RequestMessage
		responseMessage *httptest.ResponseMessage

		response = &bytes.Buffer{}
	)

	BeforeEach(func() {
		singleton = &playground.SingletonResource{Path: "/singleton"}
		response.Reset()
	})

	Describe("#Get", func() {
		Context("when no content has been written", func() {
			BeforeEach(func() {
				request = &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   "/singleton/data",
				}

				singleton.Get(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("responds 404 Not Found", func() {
				responseMessage.ShouldBeWellFormed()
				responseMessage.StatusShouldBe(404, "Not Found")
			})
			It("sets Content-Type to text/plain", func() {
				responseMessage.HeaderShould("Content-Type", Equal("text/plain"))
			})
			It("sets Content-Length to the length of the response", func() {
				responseMessage.HeaderShould("Content-Length", Equal("26"))
			})
			It("writes an error message to the message body", func() {
				responseMessage.BodyShould(Equal("Not found: /singleton/data"))
			})
		})

		Context("when content has been written", func() {
			Context("when the requested path is <path>/data", func() {
				BeforeEach(func() {
					postRequest := &httptest.RequestMessage{
						MethodReturns: http.POST,
						PathReturns: "/singleton",
					}
					postRequest.SetStringBody("field=value")
					singleton.Post(response, postRequest)

					response.Reset()
					singleton.Get(response, &httptest.RequestMessage{
						MethodReturns: http.GET,
						PathReturns:   "/singleton/data",
					})
					responseMessage = httptest.ParseResponse(response)
				})

				It("responds 200 OK", func() {
					responseMessage.ShouldBeWellFormed()
					responseMessage.StatusShouldBe(200, "OK")
				})
				It("sets Content-Type to text/plain", func() {
					responseMessage.HeaderShould("Content-Type", Equal("text/plain"))
				})
				It("sets Content-Length to the length of the response", func() {
					responseMessage.HeaderShould("Content-Length", Equal("11"))
				})
				It("writes the current data to the body", func() {
					responseMessage.BodyShould(Equal("field=value"))
				})
			})

			Context("when the requested path is anything else", func() {
				BeforeEach(func() {
					request = &httptest.RequestMessage{
						MethodReturns: http.GET,
						PathReturns:   "/singleton/missing",
					}

					singleton.Get(response, request)
					responseMessage = httptest.ParseResponse(response)
				})

				It("responds 404 Not Found", func() {
					responseMessage.StatusShouldBe(404, "Not Found")
				})
			})
		})
	})

	Describe("#Post", func() {
		Context("given any data in the body", func() {
			BeforeEach(func() {
				request = &httptest.RequestMessage{
					MethodReturns: http.POST,
					PathReturns:   "/singleton",
				}

				singleton.Post(response, request)
				responseMessage = httptest.ParseResponse(response)
			})

			It("responds 201 Created", func() {
				responseMessage.ShouldBeWellFormed()
				responseMessage.StatusShouldBe(201, "Created")
			})
			It("sets Location to <path>/data", func() {
				responseMessage.HeaderShould("Location", Equal("/singleton/data"))
			})
		})
	})
})
