package opt_test

import (
	"bufio"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/opt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Route", func() {
	var (
		router        http.Route
		controller    *ControllerMock
		requested     *http.RequestLine
		routedRequest http.Request
	)

	BeforeEach(func() {
		controller = &ControllerMock{}
		router = &opt.Route{Controller: controller}
	})

	It("routes OPTIONS * to Controller#Options", func() {
		requested = &http.RequestLine{Method: "OPTIONS", Target: "*"}
		routedRequest = router.Route(requested)
		routedRequest.Handle(&bufio.Writer{})
		controller.OptionsShouldHaveBeenCalled()
	})

	It("returns nil to pass on any other method", func() {
		requested = &http.RequestLine{Method: "GET", Target: "*"}
		routedRequest = router.Route(requested)
		Expect(routedRequest).To(BeNil())
	})

	It("returns nil to pass on any other target", func() {
		requested = &http.RequestLine{Method: "OPTIONS", Target: "/"}
		routedRequest = router.Route(requested)
		Expect(routedRequest).To(BeNil())
	})
})
