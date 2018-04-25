package teapot_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("teapotRoute", func() {
	var (
		router        http.Route
		controller    *TeapotMock
		requested     *http.RequestLine
		routedRequest http.Request
	)

	Describe("#Route", func() {
		Context("when the target is a resource that the teapot can respond to", func() {
			BeforeEach(func() {
				controller = &TeapotMock{RespondsToTarget: "/caffeine"}
				router = &teapot.Route{Resource: controller}
			})

			It("routes GET requests to that target to the teapot", func() {
				requested = &http.RequestLine{Method: "GET", Target: "/caffeine"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				controller.GetShouldHaveReceived("/caffeine")
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := &http.RequestLine{Method: "TRACE", Target: "/caffeine"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "OPTIONS")))
			})
		})

		It("passes on any other target", func() {
			controller = &TeapotMock{}
			router = &teapot.Route{Resource: controller}

			requested = &http.RequestLine{Method: "GET", Target: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
