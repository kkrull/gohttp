package fs_test

import (
	"bytes"
	"fmt"

	. "github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("DirectoryListing", func() {
	Describe("#WriteTo", func() {
		var (
			listing http.Response
			output  *bytes.Buffer
			err     error
		)

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				listing = &DirectoryListing{Files: []string{"one", "two"}}
				output = &bytes.Buffer{}
				listing.WriteTo(output)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(output.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the size of the message", func() {
				Expect(output.String()).To(containHeader("Content-Length", "8"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(output.String()).To(containHeader("Content-Type", "text/plain"))
			})
			It("lists the files in the base path", func() {
				Expect(output.String()).To(haveMessageBody("one\ntwo\n"))
			})
		})
	})
})

func haveStatus(status int, reason string) types.GomegaMatcher {
	return HavePrefix(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, reason))
}

func containHeader(name string, value string) types.GomegaMatcher {
	return ContainSubstring(fmt.Sprintf("%s: %s\r\n", name, value))
}

func haveMessageBody(message string) types.GomegaMatcher {
	return HaveSuffix(fmt.Sprintf("\r\n\r\n%s", message))
}
