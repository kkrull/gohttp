package log_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/log"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewLogRoute", func() {
	It("returns a Route at the given path", func() {
		route := log.NewLogRoute("/foo")
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&log.Route{
			Path:   "/foo",
			Viewer: &log.Viewer{},
		}))
	})
})

var _ = Describe("Route", func() {
	const (
		configuredPath = "/logs"
	)

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
