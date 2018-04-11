package http_test

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
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
					Target:        "/missing.txt"}
				err = request.Handle(response)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("Responds with 404 Not Found", func() {
				Expect(response.String()).To(haveStatus(404, "Not Found"))
			})
			It("sets Content-Length to the length of the response", func() {
				Expect(response.String()).To(containHeader("Content-Length", "23"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(response.String()).To(containHeader("Content-Type", "text/plain"))
			})
			It("writes an error message to the message body", func() {
				Expect(response.String()).To(haveMessageBody("Not found: /missing.txt"))
			})
		})

		Context("when the target is a readable file in the base path", func() {
			BeforeEach(func() {
				request = &http.GetRequest{
					BaseDirectory: basePath,
					Target:        "/readable.txt"}

				existingFile := path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
				err = request.Handle(response)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(response.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the number of bytes in the file", func() {
				Expect(response.String()).To(containHeader("Content-Length", "1"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(response.String()).To(containHeader("Content-Type", "text/plain"))
			})
			It("writes the contents of the file to the message body", func() {
				Expect(response.String()).To(haveMessageBody("A"))
			})
		})

		Context("when the target is /", func() {
			BeforeEach(func() {
				request = &http.GetRequest{
					BaseDirectory: basePath,
					Target:        "/"}

				existingFile := path.Join(basePath, "one")
				Expect(createTextFile(existingFile, "1")).To(Succeed())
				err = request.Handle(response)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(response.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the size of the message", func() {
				Expect(response.String()).To(containHeader("Content-Length", "4"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(response.String()).To(containHeader("Content-Type", "text/plain"))
			})
			It("lists the files in the base path", func() {
				Expect(response.String()).To(haveMessageBody("one\n"))
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

func createTextFile(filename string, contents string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	byteContents := bytes.NewBufferString(contents).Bytes()
	bytesWritten, err := file.Write(byteContents)
	if err != nil {
		return err
	} else if bytesWritten != len(byteContents) {
		return fmt.Errorf("expected to write %d bytes to %s, but only wrote %d", len(byteContents), filename, bytesWritten)
	}

	return nil
}

func haveStatus(status int, reason string) types.GomegaMatcher {
	return HavePrefix(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, reason))
}

func containHeader(name string, value string) types.GomegaMatcher {
	return ContainSubstring(fmt.Sprintf("%s: %s\r\n", name, value))
}

func haveMessageBody(message string) types.GomegaMatcher {
	return HaveSuffix(fmt.Sprintf("\r\n\r\n%s", message))
}
