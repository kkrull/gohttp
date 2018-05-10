package fs_test

import (
	"bytes"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DirectoryListing", func() {
	var (
		resource       *fs.DirectoryListing
		response       *httptest.ResponseMessage
		responseBuffer = &bytes.Buffer{}
	)

	BeforeEach(func() {
		responseBuffer.Reset()
	})

	Describe("#Get", func() {
		Context("always", func() {
			BeforeEach(func() {
				resource = &fs.DirectoryListing{Files: []string{}}
				resource.Get(responseBuffer, http.NewGetMessage("/"))
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds with 200 OK", func() {
				response.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to the size of the message", func() {
				contentLength, err := response.HeaderAsInt("Content-Length")
				Expect(err).NotTo(HaveOccurred())
				Expect(contentLength).To(BeNumerically(">", 0))
			})
			It("sets Content-Type to text/html", func() {
				response.HeaderShould("Content-Type", Equal("text/html"))
			})
		})

		Context("given no file names", func() {
			BeforeEach(func() {
				resource = &fs.DirectoryListing{Files: []string{}}
				resource.Get(responseBuffer, http.NewGetMessage("/"))
				response = httptest.ParseResponse(responseBuffer)
			})

			It("has an empty list of file names", func() {
				response.BodyShould(MatchRegexp(".*<ul>\\s*<[/]ul>.*"))
			})
		})

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				resource = &fs.DirectoryListing{
					Files:      []string{"one", "two"},
					HrefPrefix: "/files"}
				resource.Get(responseBuffer, http.NewGetMessage("/"))
				response = httptest.ParseResponse(responseBuffer)
			})

			It("lists links to the files, using absolute paths with the given prefix", func() {
				response.BodyShould(ContainSubstring("<a href=\"/files/one\">one</a>"))
				response.BodyShould(ContainSubstring("<a href=\"/files/two\">two</a>"))
			})
		})
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			resource = &fs.DirectoryListing{Files: []string{}}
			resource.Head(responseBuffer, http.NewHeadMessage("/"))
			response = httptest.ParseResponse(responseBuffer)
		})

		It("returns the same status as #Get", func() {
			response.StatusShouldBe(200, "OK")
		})
		It("sets the same headers as #Get", func() {
			contentLength, _ := response.HeaderAsInt("Content-Length")
			Expect(contentLength).To(BeNumerically(">", 0))
			response.HeaderShould("Content-Type", HavePrefix("text/html"))
		})
		It("has no body", func() {
			response.BodyShould(BeEmpty())
		})
	})

	Describe("#Options", func() {
		BeforeEach(func() {
			resource = &fs.DirectoryListing{}
			requestMessage := http.NewOptionsMessage("/")
			request := requestMessage.MakeResourceRequest(resource)

			request.Handle(responseBuffer)
		})

		It("supports read operations", httptest.ShouldAllowMethods(responseBuffer, "GET", "HEAD", "OPTIONS"))
	})
})
