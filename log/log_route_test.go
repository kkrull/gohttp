package log_test

import (
	"bytes"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/log"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const configuredPath = "/logs"

var _ = Describe("::NewLogRoute", func() {
	It("returns a Route at the given path", func() {
		logger := &RequestBufferStub{}
		route := log.NewLogRoute("/foo", logger)
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&log.Route{
			Path: "/foo",
			Viewer: &log.Viewer{
				Requests: logger,
			},
		}))
	})
})

var _ = Describe("Route", func() {
	var (
		router   http.Route
		response = &bytes.Buffer{}
	)

	Describe("#Route", func() {
		BeforeEach(func() {
			router = &log.Route{Path: configuredPath}
			response.Reset()
		})

		Context("when the path is the given path", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage(configuredPath)
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows OPTIONS", httptest.ShouldAllowMethods(response, http.OPTIONS))
				It("allows GET", httptest.ShouldAllowMethods(response, http.GET))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage(configuredPath)
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(clienterror.MethodNotAllowed()))
			})
		})

		It("passes on any other path by returning nil", func() {
			requested := http.NewGetMessage("/")
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})

var _ = Describe("Viewer", func() {
	var (
		viewer   *log.Viewer
		logger   log.RequestBuffer
		response *httptest.ResponseMessage
	)

	BeforeEach(func() {
		logger = &RequestBufferStub{
			NumBytesReturns: 4,
			WriteToWill:     bytes.NewBufferString("ASDF"),
		}
		viewer = &log.Viewer{Requests: logger}
	})

	Describe("#Get", func() {
		const (
			validAuthorization           = "Basic YWRtaW46aHVudGVyMg=="
			invalidAuthorizationMethod   = "Wat YWRtaW46aHVudGVyMg=="
			invalidAuthorizationPassword = "Basic YWRtaW46dG90YWxseWJvZ3VzcGFzc3dvcmQK"
		)

		Context("given no Authorization header", func() {
			BeforeEach(func() {
				response = invokeResourceMethod(viewer.Get, http.NewGetMessage(configuredPath))
			})

			It("responds 401 Unauthorized", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(401, "Unauthorized")
			})
			It("sets WWW-Authenticate to a Basic challenge in the logs realm", func() {
				response.HeaderShould("WWW-Authenticate", Equal("Basic realm=\"logs\""))
			})
		})

		Context("given an Authorization header with valid Basic credentials", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   configuredPath,
					TargetReturns: configuredPath,
				}
				request.AddHeader("Authorization", validAuthorization)
				response = invokeResourceMethod(viewer.Get, request)
			})

			It("responds 200 OK", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(200, "OK")
			})
			It("responds with the contents of the configured readable buffer", func() {
				response.HeaderShould("Content-Length", Equal("4"))
				response.BodyShould(Equal("ASDF"))
			})
		})

		Context("given an Authorization header with any other method", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   configuredPath,
					TargetReturns: configuredPath,
				}
				request.AddHeader("Authorization", invalidAuthorizationMethod)
				response = invokeResourceMethod(viewer.Get, request)
			})

			It("responds 403 Forbidden", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(403, "Forbidden")
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})

		Context("given an Authorization header with invalid Basic credentials", func() {
			BeforeEach(func() {
				request := &httptest.RequestMessage{
					MethodReturns: http.GET,
					PathReturns:   configuredPath,
					TargetReturns: configuredPath,
				}
				request.AddHeader("Authorization", invalidAuthorizationPassword)
				response = invokeResourceMethod(viewer.Get, request)
			})

			It("responds 403 Forbidden", func() {
				response.ShouldBeWellFormed()
				response.StatusShouldBe(403, "Forbidden")
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})
	})
})

func invokeResourceMethod(invokeMethod httpResourceMethod, request http.RequestMessage) *httptest.ResponseMessage {
	response := &bytes.Buffer{}
	invokeMethod(response, request)
	return httptest.ParseResponse(response)
}

type httpResourceMethod = func(io.Writer, http.RequestMessage)
