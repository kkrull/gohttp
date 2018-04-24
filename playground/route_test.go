package playground_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("::NewRoute", func() {
	It("returns a Route configured for this package", func() {
		route := playground.NewRoute()
		Expect(route).NotTo(BeNil())
		Expect(route.Writable).To(BeAssignableToTypeOf(&playground.ReadWriteNopResource{}))
	})
})

var _ = Describe("Route", func() {
	Describe("#Route", func() {
		var (
			router            http.Route
			readOnlyResource  *ReadOnlyResourceMock
			readWriteResource *ReadWriteResourceMock
		)

		BeforeEach(func() {
			readOnlyResource = &ReadOnlyResourceMock{}
			readWriteResource = &ReadWriteResourceMock{}
			router = &playground.Route{
				Readable: readOnlyResource,
				Writable: readWriteResource,
			}
		})

		Context("when the target is /method_options", func() {
			It("routes GET to Writable#Get", func() {
				handleRequest(router, "GET", "/method_options")
				readWriteResource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Writable#Head", func() {
				handleRequest(router, "HEAD", "/method_options")
				readWriteResource.HeadShouldHaveBeenCalled()
			})

			It("routes OPTIONS to Writable#Options", func() {
				handleRequest(router, "OPTIONS", "/method_options")
				readWriteResource.OptionsShouldHaveBeenCalled()
			})

			It("routes POST to Writable#Post", func() {
				handleRequest(router, "POST", "/method_options")
				readWriteResource.PostShouldHaveBeenCalled()
			})

			It("routes PUT to Writable#Put", func() {
				handleRequest(router, "PUT", "/method_options")
				readWriteResource.PutShouldHaveBeenCalled()
			})
		})

		Context("when the target is /method_options2", func() {
			It("routes GET to Writable#Get", func() {
				handleRequest(router, "GET", "/method_options2")
				readOnlyResource.GetShouldHaveBeenCalled()
			})

			It("routes HEAD to Writable#Options", func() {
				handleRequest(router, "HEAD", "/method_options2")
				readOnlyResource.HeadShouldHaveBeenCalled()
			})

			It("routes OPTIONS to Writable#Options", func() {
				handleRequest(router, "OPTIONS", "/method_options2")
				readOnlyResource.OptionsShouldHaveBeenCalled()
			})

			It("returns nil for any other method", func() {
				requested := &http.RequestLine{Method: "PUT", Target: "/method_options2"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeNil())
			})
		})

		It("returns nil on any other target", func() {
			requested := &http.RequestLine{Method: "GET", Target: "/"}
			routedRequest := router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})

func handleRequest(router http.Route, method, target string) {
	requested := &http.RequestLine{Method: method, Target: target}
	routedRequest := router.Route(requested)
	ExpectWithOffset(1, routedRequest).NotTo(BeNil())
	routedRequest.Handle(&bufio.Writer{})
}
