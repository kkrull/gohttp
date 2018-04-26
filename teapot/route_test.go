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
		teapotMock    *TeapotMock
		requested     *http.RequestLine
		routedRequest http.Request
	)

	Describe("#Route", func() {
		Context("when the target is a resource that the teapot can respond to", func() {
			BeforeEach(func() {
				teapotMock = &TeapotMock{RespondsToTarget: "/caffeine"}
				router = &teapot.Route{Teapot: teapotMock}
			})

			It("routes GET requests to that target to the teapot", func() {
				requested = &http.RequestLine{TheMethod: "GET", TheTarget: "/caffeine"}
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				teapotMock.GetShouldHaveReceived("/caffeine")
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := &http.RequestLine{TheMethod: "TRACE", TheTarget: "/caffeine"}
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "OPTIONS")))
			})
		})

		It("passes on any other target", func() {
			teapotMock = &TeapotMock{}
			router = &teapot.Route{Teapot: teapotMock}

			requested = &http.RequestLine{TheMethod: "GET", TheTarget: "/file.txt"}
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
