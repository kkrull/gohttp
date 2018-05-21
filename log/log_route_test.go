package log_test

import (
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewLogRoute", func() {
	It("returns a Route at the given path", func() {
		route := log.NewLogRoute("/foo")
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&log.Route{}))
	})
})

var _ = Describe("Route", func() {
	var (
		router http.Route
	)

	Describe("#Route", func() {
		BeforeEach(func() {
			router = &log.Route{}
		})

		It("passes on any other path by returning nil", func() {
			requested := http.NewGetMessage("/")
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})
