package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewCookieRoute", func() {
	It("returns a Route at the given paths", func() {
		route := playground.NewCookieRoute("/set", "/read")
		Expect(route).NotTo(BeNil())
		Expect(route).To(BeEquivalentTo(&playground.CookieRoute{
			SetTypePath: "/set",
			ReadTypePath: "/read",
		}))
	})
})

var _ = Describe("CookieRoute", func() {
	var (
		router   http.Route
		response = &bytes.Buffer{}
	)

	Describe("#Route", func() {
		BeforeEach(func() {
			router = &playground.CookieRoute{
			}
			response.Reset()
		})

		It("passes on any other path by returning nil", func() {
			requested := http.NewGetMessage("/")
			Expect(router.Route(requested)).To(BeNil())
		})
	})
})
