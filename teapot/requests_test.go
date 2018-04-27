package teapot_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/teapot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IdentityTeapot", func() {
	var (
		theTeapot      teapot.Teapot
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#RespondsTo", func() {
		It("returns true for /coffee", func() {
			theTeapot = &teapot.IdentityTeapot{}
			Expect(theTeapot.RespondsTo("/coffee")).To(BeTrue())
		})

		It("returns true for /tea", func() {
			theTeapot = &teapot.IdentityTeapot{}
			Expect(theTeapot.RespondsTo("/tea")).To(BeTrue())
		})
	})

	Describe("#Get", func() {
		Context("when the target is /coffee", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				theTeapot = &teapot.IdentityTeapot{}

				theTeapot.Get(responseBuffer, http.NewRequestMessage("GET", "/coffee"))
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

		Context("when the target is /tea", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				theTeapot = &teapot.IdentityTeapot{}

				theTeapot.Get(responseBuffer, http.NewRequestMessage("GET", "/tea"))
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
})
