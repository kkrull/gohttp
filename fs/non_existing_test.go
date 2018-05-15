package fs_test

import (
	"bytes"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NonExisting", func() {
	var (
		resource       *fs.NonExisting
		response       *httptest.ResponseMessage
		responseBuffer = &bytes.Buffer{}
	)

	BeforeEach(func() {
		responseBuffer.Reset()
	})

	Describe("#Get", func() {
		BeforeEach(func() {
			resource = &fs.NonExisting{Path: "/missing.txt"}
			resource.Get(responseBuffer, http.NewGetMessage("/missing.txt"))
			response = httptest.ParseResponse(responseBuffer)
		})

		It("Responds with 404 Not Found", func() {
			response.StatusShouldBe(404, "Not Found")
		})
		It("sets Content-Length to the length of the response", func() {
			response.HeaderShould("Content-Length", Equal("23"))
		})
		It("sets Content-Type to text/plain", func() {
			response.HeaderShould("Content-Type", Equal("text/plain"))
		})
		It("writes an error message to the message body", func() {
			response.BodyShould(Equal("Not found: /missing.txt"))
		})
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			resource = &fs.NonExisting{Path: "/missing.txt"}
			resource.Head(responseBuffer, http.NewHeadMessage("/missing.txt"))
			response = httptest.ParseResponse(responseBuffer)
		})

		It("Responds with 404 Not Found", func() {
			response.StatusShouldBe(404, "Not Found")
		})
		It("sets Content-Length to the length of the response", func() {
			response.HeaderShould("Content-Length", Equal("23"))
		})
		It("sets Content-Type to text/plain", func() {
			response.HeaderShould("Content-Type", Equal("text/plain"))
		})
		It("has no body", func() {
			response.BodyShould(BeEmpty())
		})
	})

	Describe("#Options", func() {
		BeforeEach(func() {
			resource = &fs.NonExisting{Path: "/missing.txt"}
			requestMessage := http.NewOptionsMessage("/missing.txt")
			request := requestMessage.MakeResourceRequest(resource)

			request.Handle(responseBuffer)
		})

		It("supports read operations", httptest.ShouldAllowMethods(responseBuffer, "GET", "HEAD", "OPTIONS"))
	})
})
