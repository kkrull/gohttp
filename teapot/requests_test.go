package teapot_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("IdentityController", func() {
	var (
		controller     teapot.Controller
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Get", func() {
		Context("given /coffee", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &teapot.IdentityController{}

				controller.Get(responseBuffer, "/coffee")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("writes a valid HTTP response", func() {
				response.ShouldBeWellFormed()
			})
			It("responds with 418 I'm a teapot", func() {
				response.StatusShouldBe(418, "I'm a teapot")
			})
		})
	})
})
