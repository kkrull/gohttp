package playground_test

import (
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRoute", func() {
	It("returns the Route for this package", func() {
		route := playground.NewRoute()
		Expect(route).NotTo(BeNil())
	})
})

var _ = Describe("Route", func() {
	Describe("#Route", func() {
		var (
			router        http.Route
			requested     *http.RequestLine
			routedRequest http.Request
		)

		BeforeEach(func() {
			router = &playground.Route{}
		})

		It("returns nil on any other request", func() {
			requested = &http.RequestLine{Method: "GET", Target: "/"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
