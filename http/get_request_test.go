package http_test

import (
	"bufio"
	"bytes"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRequest", func() {
	Describe("#Handle", func() {
		var (
			request  *http.GetRequest
			basePath string
			response *bytes.Buffer
			err      error
		)

		BeforeEach(func() {
			basePath = makeEmptyTestDirectory("GetRequest", os.ModePerm)
			response = &bytes.Buffer{}
		})

		Context("when the resolved target does not exist", func() {
			BeforeEach(func() {
				request = &http.GetRequest{
					BaseDirectory: basePath,
					Target:        "/missing.txt",
					Version:       "HTTP/1.1"}
				err = request.Handle(bufio.NewWriter(response))
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("Responds with 404 Not Found", func() {
				Expect(response.String()).To(Equal("HTTP/1.1 404 Not Found\r\n"))
			})
		})
	})
})

func makeEmptyTestDirectory(testName string, fileMode os.FileMode) string {
	testPath := path.Join(".test", testName)
	Expect(os.RemoveAll(testPath)).To(Succeed())
	Expect(os.MkdirAll(testPath, fileMode)).To(Succeed())
	return testPath
}
