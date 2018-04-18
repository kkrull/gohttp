package fs_test

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	var (
		controller *fs.Controller
		basePath   string
		response   *bytes.Buffer
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("Controller", os.ModePerm)
		controller = &fs.Controller{BaseDirectory: basePath}
		response = &bytes.Buffer{}
	})

	Describe("#Get", func() {
		Context("when the resolved target does not exist", func() {
			BeforeEach(func() {
				controller.Get(response, "/missing.txt")
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

		Context("when the target is a readable text file in the base path", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
				controller.Get(response, "/readable.txt")
			})

			It("responds with 200 OK", func() {
				Expect(response.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the number of bytes in the file", func() {
				Expect(response.String()).To(containHeader("Content-Length", "1"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(response.String()).To(containHeader("Content-Type", "text/plain; charset=utf-8"))
			})
			It("writes the contents of the file to the message body", func() {
				Expect(response.String()).To(haveMessageBody("A"))
			})
		})

		Context("when the target is a readable file named with a registered extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "image.jpeg")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to the MIME type registered for that extension", func() {
				controller.Get(response, "/image.jpeg")
				Expect(response.String()).To(containHeader("Content-Type", "image/jpeg"))
			})
		})

		Context("when the target is a readable file without an extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "assumed-text")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to text/plain", func() {
				controller.Get(response, "/assumed-text")
				Expect(response.String()).To(containHeader("Content-Type", "text/plain"))
			})
		})

		Context("when the target is /", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "one")
				Expect(createTextFile(existingFile, "1")).To(Succeed())
			})

			It("responds with 200 OK", func() {
				controller.Get(response, "/")
				Expect(response.String()).To(haveStatus(200, "OK"))
			})
		})
	})

	Describe("#Head", func() {
		var getResponse = &bytes.Buffer{}

		BeforeEach(func() {
			controller.Get(getResponse, "/missing.txt")
			controller.Head(response, "/missing.txt")
		})

		It("responds with the same status as #Get would have", func() {
			Expect(getResponse.String()).To(haveStatus(404, "Not Found"))
			Expect(response.String()).To(haveStatus(404, "Not Found"))
		})
		It("responds with the same headers as #Get would have", func() {
			Expect(response.String()).To(containHeader("Content-Type", "text/plain"))
			Expect(response.String()).To(containHeader("Content-Length", "23"))
		})
		It("does not write a message body", func() {
			Expect(response.String()).To(HaveSuffix("\r\n\r\n"))
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
