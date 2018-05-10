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
	XIt("DirectoryListing implements Get and Head")

	Describe("#WriteTo", func() {
		var (
			listing http.Response
			message *httptest.ResponseMessage
			err     error
		)

		Context("always", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{Files: []string{}}
				listing.WriteTo(output)
				message = httptest.ParseResponse(output)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				message.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to the size of the message", func() {
				contentLength, err := message.HeaderAsInt("Content-Length")
				Expect(err).NotTo(HaveOccurred())
				Expect(contentLength).To(BeNumerically(">", 0))
			})
			It("sets Content-Type to text/html", func() {
				message.HeaderShould("Content-Type", Equal("text/html"))
			})
		})

		Context("given no file names", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{Files: []string{}}
				listing.WriteTo(output)
				message = httptest.ParseResponse(output)
			})

			It("has an empty list of file names", func() {
				message.BodyShould(MatchRegexp(".*<ul>\\s*<[/]ul>.*"))
			})
		})

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{
					Files:      []string{"one", "two"},
					HrefPrefix: "/files"}
				listing.WriteTo(output)
				message = httptest.ParseResponse(output)
			})

			It("lists links to the files, using absolute paths with the given prefix", func() {
				message.BodyShould(ContainSubstring("<a href=\"/files/one\">one</a>"))
				message.BodyShould(ContainSubstring("<a href=\"/files/two\">two</a>"))
			})
		})
	})
})
