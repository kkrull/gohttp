package playground_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/playground"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("ReadableNopResource", func() {
	var (
		controller *playground.ReadableNopResource
		response   = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
		controller = &playground.ReadableNopResource{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Options", func() {
		BeforeEach(func() {
			controller.Options(response)
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

var _ = Describe("ReadWriteNopResource", func() {
	var (
		controller *playground.ReadWriteNopResource
		response   = &bytes.Buffer{}
	)

	BeforeEach(func() {
		response.Reset()
		controller = &playground.ReadWriteNopResource{}
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			controller.Get(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Options", func() {
		BeforeEach(func() {
			controller.Options(response)
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

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(response)
		})

		It("responds 200 OK with no body", ShouldHaveNoBody(response, 200, "OK"))
	})
})
