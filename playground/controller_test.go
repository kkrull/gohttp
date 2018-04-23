package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("WritableNopController", func() {
	var (
		controller *playground.WritableNopController
		response   = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
		controller = &playground.WritableNopController{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(response, "/method_options")
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response, "/method_options")
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Options", func() {
		Context("given a request for /method_options", func() {
			BeforeEach(func() {
				controller.Options(response, "/method_options")
			})

			It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				responseMessage := httptest.ParseResponse(response)
				responseMessage.HeaderShould("Allow", ContainSubstrings([]string{
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
				controller.Options(response, "/method_options2")
			})

			It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
			It("sets Allow to the methods that SimpleOption expects for this route", func() {
				responseMessage := httptest.ParseResponse(response)
				responseMessage.HeaderShould("Allow", ContainSubstrings([]string{
					"GET",
					"HEAD",
					"OPTIONS",
				}))
			})
		})
	})

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(response, "/method_options")
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(response, "/method_options")
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})
})

func ShouldHaveNoBody(response *bytes.Buffer, status int, reason string) func() {
	return func() {
		responseMessage := httptest.ParseResponse(response)
		responseMessage.ShouldBeWellFormed()
		responseMessage.StatusShouldBe(status, reason)
		responseMessage.HeaderShould("Content-Length", Equal("0"))
		responseMessage.BodyShould(BeEmpty())
	}
}

func ContainSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}
