package http_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
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
				Expect(response.String()).To(HavePrefix("HTTP/1.1 404 Not Found\r\n"))
			})
			It("sets Content-Length to 0", func() {
				Expect(response.String()).To(ContainSubstring("Content-Length: 0\r\n"))
			})
			It("has no body", func() {
				Expect(response.String()).To(HaveSuffix("\r\n\r\n"))
			})
		})

		Context("when the target is a readable file in the specified path", func() {
			BeforeEach(func() {
				request = &http.GetRequest{
					BaseDirectory: basePath,
					Target:        "/readable.txt",
					Version:       "HTTP/1.1"}

				existingFile := path.Join(basePath, "readable.txt")

				file, err := os.Create(existingFile)
				Expect(err).NotTo(HaveOccurred())
				contents := bytes.NewBufferString("A")
				file.Write(contents.Bytes())

				Expect(ioutil.ReadFile(existingFile)).NotTo(BeEmpty())
				err = request.Handle(bufio.NewWriter(response))
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(response.String()).To(HavePrefix("HTTP/1.1 200 OK\r\n"))
			})
			It("sets Content-Length to the number of bytes in the file", func() {
				Expect(response.String()).To(ContainSubstring("Content-Length: 1\r\n"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(response.String()).To(ContainSubstring("Content-Type: text/plain\r\n"))
			})
			It("writes the contents of the file to the message body", func() {
				Expect(response.String()).To(HaveSuffix("\r\n\r\nA"))
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
