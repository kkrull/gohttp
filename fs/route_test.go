package fs_test

import (
	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("route", func() {
	var (
		router = fs.NewRoute("/public")
	)

	Describe("#Route", func() {
		It("routes GET requests to GetRequest", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/foo"}
			Expect(router.Route(requested)).To(BeEquivalentTo(
				&fs.GetRequest{
					BaseDirectory: "/public",
					Target:        "/foo",
				}))
		})
	})
})
