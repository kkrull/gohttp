package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("AllowedMethodsController", func() {
	var (
		controller     *playground.AllowedMethodsController
		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	ShouldHaveNoBody := func(status int, reason string) func() {
		return func() {
			response.ShouldBeWellFormed()
			response.StatusShouldBe(status, reason)
			response.HeaderShould("Content-Length", Equal("0"))
			response.BodyShould(BeEmpty())
		}
	}

	BeforeEach(func() {
		responseBuffer = &bytes.Buffer{}
		controller = &playground.AllowedMethodsController{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(responseBuffer, "/method_options")
			response = httptest.ParseResponse(responseBuffer)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(responseBuffer, "/method_options")
			response = httptest.ParseResponse(responseBuffer)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
	})

	Describe("#Options", func() {
		Context("given a request for /method_options", func() {
			BeforeEach(func() {
				controller.Options(responseBuffer, "/method_options")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				response.HeaderShould("Allow", ContainSubstrings([]string{
					"GET",
					"HEAD",
					"OPTIONS",
					"POST",
					"PUT",
				}))
			})
		})

		Context("given a request for /method_options2", func() {
			BeforeEach(func() {
				controller.Options(responseBuffer, "/method_options2")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				response.HeaderShould("Allow", ContainSubstrings([]string{
					"GET",
					"HEAD",
					"OPTIONS",
				}))
			})
		})
	})

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(responseBuffer, "/method_options")
			response = httptest.ParseResponse(responseBuffer)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(responseBuffer, "/method_options")
			response = httptest.ParseResponse(responseBuffer)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(200, "OK"))
	})
})

func ContainSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}
