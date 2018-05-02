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
		requested     http.RequestMessage
		routedRequest http.Request
	)

	Describe("#Route", func() {
		Context("when the path is a resource that the teapot can respond to", func() {
			BeforeEach(func() {
				teapotMock = &TeapotMock{RespondsToPath: "/caffeine"}
				router = &teapot.Route{Teapot: teapotMock}
			})

			It("routes GET requests to that path to the teapot", func() {
				requested = http.NewGetMessage("/caffeine")
				routedRequest = router.Route(requested)
				routedRequest.Handle(&bufio.Writer{})
				teapotMock.GetShouldHaveReceived("/caffeine")
			})

			It("returns MethodNotAllowed for any other method", func() {
				requested := http.NewTraceMessage("/caffeine")
				routedRequest := router.Route(requested)
				Expect(routedRequest).To(BeEquivalentTo(clienterror.MethodNotAllowed("GET", "OPTIONS")))
			})
		})

		It("passes on any other path", func() {
			teapotMock = &TeapotMock{}
			router = &teapot.Route{Teapot: teapotMock}

			requested = http.NewGetMessage("/file.txt")
			routedRequest = router.Route(requested)
			Expect(routedRequest).To(BeNil())
		})
	})
})
