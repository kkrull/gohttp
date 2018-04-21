package teapot_test

import (
	"bytes"

	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("IdentityController", func() {
	var (
		controller teapot.Controller
		//response       *HttpMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Get", func() {
		Context("given /coffee", func() {
			XIt("responds 418 I'm a teapot", func() {
				controller = &teapot.IdentityController{}
				controller.Get(responseBuffer, "/coffee")
				//response.StatusShouldBe(418, "I'm a teapot")
			})
		})
	})
})
