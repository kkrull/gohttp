package coffee_test

import (
	"github.com/kkrull/gohttp/coffee"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("coffeeRoute", func() {
	var (
		router        = coffee.NewRoute()
		requested     *http.RequestLine
		routedRequest http.Request
	)

	Describe("#Route", func() {
		It("routes GET /coffee to GetCoffeeRequest", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/coffee"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeAssignableToTypeOf(&coffee.GetCoffeeRequest{}))
		})

		It("passes on any other target", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
