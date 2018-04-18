package fs_test

import (
	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("route", func() {
	var (
		router     = fs.NewRoute("/public")
		controller = &fs.Controller{BaseDirectory: "/public"}
	)

	Describe("#Route", func() {
		It("routes GET requests to GetRequest", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/foo"}
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeEquivalentTo(
				&fs.GetRequest{
					Controller: controller,
					Target:     "/foo",
				}))
		})

		It("routes HEAD requests to HeadRequest", func() {
			requested := &http.RequestLine{Method: "HEAD", Target: "/foo"}
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeEquivalentTo(
				&fs.HeadRequest{
					Controller: controller,
					Target:     "/foo",
				}))
		})

		It("passes on any other method by returning nil", func() {
			requested := &http.RequestLine{Method: "TRACE", Target: "/foo"}
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
