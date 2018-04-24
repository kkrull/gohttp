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

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response)
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
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

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			controller.Head(response)
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Post", func() {
		BeforeEach(func() {
			controller.Post(response)
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})

	Describe("#Put", func() {
		BeforeEach(func() {
			controller.Put(response)
		})

		It("responds 200 OK with no body", httptest.ShouldHaveNoBody(response, 200, "OK"))
	})
})
