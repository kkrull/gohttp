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
		controller     *fs.Controller
		basePath       string
		responseBuffer *bytes.Buffer
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("Controller", os.ModePerm)
		controller = &fs.Controller{BaseDirectory: basePath}
		responseBuffer = &bytes.Buffer{}
	})

	Describe("#Get", func() {
		Context("when the resolved target does not exist", func() {
			BeforeEach(func() {
				controller.Get(responseBuffer, "/missing.txt")
			})

			It("Responds with 404 Not Found", func() {
				Expect(responseBuffer.String()).To(haveStatus(404, "Not Found"))
			})
			It("sets Content-Length to the length of the response", func() {
				Expect(responseBuffer.String()).To(containHeader("Content-Length", "23"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(responseBuffer.String()).To(containHeader("Content-Type", "text/plain"))
			})
			It("writes an error message to the message body", func() {
				Expect(responseBuffer.String()).To(haveMessageBody("Not found: /missing.txt"))
			})
		})

		Context("when the target is a readable text file in the base path", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
				controller.Get(responseBuffer, "/readable.txt")
			})

			It("responds with 200 OK", func() {
				Expect(responseBuffer.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the number of bytes in the file", func() {
				Expect(responseBuffer.String()).To(containHeader("Content-Length", "1"))
			})
			It("sets Content-Type to text/plain", func() {
				Expect(responseBuffer.String()).To(containHeader("Content-Type", "text/plain; charset=utf-8"))
			})
			It("writes the contents of the file to the message body", func() {
				Expect(responseBuffer.String()).To(haveMessageBody("A"))
			})
		})

		Context("when the target is a readable file named with a registered extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "image.jpeg")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to the MIME type registered for that extension", func() {
				controller.Get(responseBuffer, "/image.jpeg")
				Expect(responseBuffer.String()).To(containHeader("Content-Type", "image/jpeg"))
			})
		})

		Context("when the target is a readable file without an extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "assumed-text")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to text/plain", func() {
				controller.Get(responseBuffer, "/assumed-text")
				Expect(responseBuffer.String()).To(containHeader("Content-Type", "text/plain"))
			})
		})

		Context("when the target is /", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "one")
				Expect(createTextFile(existingFile, "1")).To(Succeed())
			})

			It("responds with 200 OK", func() {
				controller.Get(responseBuffer, "/")
				Expect(responseBuffer.String()).To(haveStatus(200, "OK"))
			})
		})
	})

	Describe("#Head", func() {
		var (
			getResponseBuffer = &bytes.Buffer{}
			response          *HttpMessage
		)

		BeforeEach(func() {
			controller.Get(getResponseBuffer, "/missing.txt")
			controller.Head(responseBuffer, "/missing.txt")
			response = &HttpMessage{Text: responseBuffer.String()}
		})

		It("responds with the same status as #Get would have", func() {
			Expect(getResponseBuffer.String()).To(haveStatus(404, "Not Found"))
			response.StatusShouldBe(404, "Not Found")
		})
		It("responds with the same headers as #Get would have", func() {
			response.ShouldHaveHeader("Content-Type", Equal("text/plain"))
			response.ShouldHaveHeader("Content-Length", Equal("23"))
		})
		It("does not write a response body", func() {
			response.BodyShould(BeEmpty())
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
