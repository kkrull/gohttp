package teapot_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IdentityController", func() {
	var (
		controller     teapot.Controller
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#GetCoffee", func() {
		BeforeEach(func() {
			responseBuffer = &bytes.Buffer{}
			controller = &teapot.IdentityController{}

			controller.GetCoffee(responseBuffer)
			response = httptest.ParseResponse(responseBuffer)
		})

		It("writes a valid HTTP response", func() {
			response.ShouldBeWellFormed()
		})
		It("sets status to 418 I'm a teapot", func() {
			response.StatusShouldBe(418, "I'm a teapot")
		})
		It("sets Content-Length to the length of the body", func() {
			response.HeaderShould("Content-Length", Equal("12"))
		})
		It("sets Content-Type to text/plain", func() {
			response.HeaderShould("Content-Type", HavePrefix("text/plain"))
		})
		It("writes I'm a teapot to the body", func() {
			response.BodyShould(Equal("I'm a teapot"))
		})

	})

	Describe("#GetTea", func() {
		BeforeEach(func() {
			responseBuffer = &bytes.Buffer{}
			controller = &teapot.IdentityController{}

			controller.GetTea(responseBuffer)
			response = httptest.ParseResponse(responseBuffer)
		})

		It("writes a valid HTTP response", func() {
			response.ShouldBeWellFormed()
		})
		It("sets status to 200 OK", func() {
			response.StatusShouldBe(200, "OK")
		})
		It("sets Content-Length to the length of the body", func() {
			response.HeaderShould("Content-Length", Equal("0"))
		})
		It("has no body", func() {
			response.BodyShould(BeEmpty())
		})
	})
})
