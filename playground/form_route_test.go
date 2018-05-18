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

var _ = Describe("::NewFormRoute", func() {
	It("returns a FormRoute", func() {
		route := playground.NewFormRoute()
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&playground.FormRoute{
			Form: &playground.SingletonForm{},
		}))
	})
})

var _ = Describe("FormRoute", func() {
	Describe("#Route", func() {
		var (
			router   http.Route
			response = &bytes.Buffer{}
		)

		BeforeEach(func() {
			router = &playground.FormRoute{}
			response.Reset()
		})

		Context("when the path is /form", func() {
			Context("when the method is OPTIONS", func() {
				BeforeEach(func() {
					requested := http.NewOptionsMessage("/form")
					routedRequest := router.Route(requested)
					Expect(routedRequest).NotTo(BeNil())
					routedRequest.Handle(response)
				})

				It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
				It("allows POST", httptest.ShouldAllowMethods(response, http.POST, http.OPTIONS))
			})

			It("replies Method Not Allowed on any other method", func() {
				requested := http.NewTraceMessage("/form")
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
