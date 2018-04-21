package opt_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/opt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaticCapabilityController", func() {
	var (
		controller     opt.ServerCapabilityController
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	BeforeEach(func() {
		responseBuffer = &bytes.Buffer{}
		controller = &opt.StaticCapabilityController{}
		controller.Options(responseBuffer)
		response = httptest.ParseResponse(responseBuffer)
	})

	Describe("#Options", func() {
		It("writes a well-formed response", func() {
			response.ShouldBeWellFormed()
		})
		It("sets status to 200 OK", func() {
			response.StatusShouldBe(200, "OK")
		})
		XIt("sets Allow")
		It("sets Content-Length to 0", func() {
			response.HeaderShould("Content-Length", Equal("0"))
		})
	})
})
