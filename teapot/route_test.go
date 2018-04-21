package teapot_test

import (
	"github.com/kkrull/gohttp/teapot"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("teapotRoute", func() {
	var (
		router        = teapot.NewRoute()
		requested     *http.RequestLine
		routedRequest http.Request
	)

	Describe("#Route", func() {
		Context("given GET /coffee", func() {
			It("routes to GetCoffeeRequest", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/coffee"}
				routedRequest = router.Route(requested)
				Expect(routedRequest).To(BeAssignableToTypeOf(&teapot.GetCoffeeRequest{}))
			})

			XIt("calls Controller#GetCoffee", func() {

			})
		})

		It("passes on any other target", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
