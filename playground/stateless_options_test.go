package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("StatelessOptions", func() {
	var (
		controller     *playground.StatelessOptionController
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Options", func() {
		BeforeEach(func() {
			responseBuffer = &bytes.Buffer{}
			controller = &playground.StatelessOptionController{}
		})

		It("responds with the 5 methods in cob_spec", func() {
			controller.Options(responseBuffer, "/anywhere")
			response = httptest.ParseResponse(responseBuffer)
		})
	})
})
