package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("StatelessOptions", func() {
	var (
		controller     *playground.StatelessOptionController
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	Describe("#Options", func() {
		Context("given a request for /method_options", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &playground.StatelessOptionController{}

				controller.Options(responseBuffer, "/method_options")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("writes a well-formed response", func() {
				response.ShouldBeWellFormed()
			})
			It("sets status to 200 OK", func() {
				response.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to 0", func() {
				response.HeaderShould("Content-Length", Equal("0"))
			})
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				response.HeaderShould("Allow", ContainSubstrings([]string{
					"GET",
					"HEAD",
					"OPTIONS",
					"POST",
					"PUT",
				}))
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})

		Context("given a request for /method_options2", func() {
			BeforeEach(func() {
				responseBuffer = &bytes.Buffer{}
				controller = &playground.StatelessOptionController{}

				controller.Options(responseBuffer, "/method_options2")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("writes a well-formed response", func() {
				response.ShouldBeWellFormed()
			})
			It("sets status to 200 OK", func() {
				response.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to 0", func() {
				response.HeaderShould("Content-Length", Equal("0"))
			})
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				response.HeaderShould("Allow", ContainSubstrings([]string{
					"GET",
					"HEAD",
					"OPTIONS",
				}))
			})
			It("has no body", func() {
				response.BodyShould(BeEmpty())
			})
		})
	})
})

func ContainSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}
