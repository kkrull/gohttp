package fs_test

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	var (
		controller *fs.Controller
		basePath   string

		response       *httptest.ResponseMessage
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

		Context("when the target is a readable text file in the base path", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
				controller.Get(responseBuffer, "/readable.txt")
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds with 200 OK", func() {
				response.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to the number of bytes in the file", func() {
				response.HeaderShould("Content-Length", Equal("1"))
			})
			It("sets Content-Type to text/plain", func() {
				response.HeaderShould("Content-Type", HavePrefix("text/plain"))
			})
			It("writes the contents of the file to the message body", func() {
				response.BodyShould(Equal("A"))
			})
		})

		Context("when the target is a readable file named with a registered extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "image.jpeg")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to the MIME type registered for that extension", func() {
				controller.Get(responseBuffer, "/image.jpeg")
				response = httptest.ParseResponse(responseBuffer)
				response.HeaderShould("Content-Type", Equal("image/jpeg"))
			})
		})

		Context("when the target is a readable file without an extension", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "assumed-text")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to text/plain", func() {
				controller.Get(responseBuffer, "/assumed-text")
				response = httptest.ParseResponse(responseBuffer)
				response.HeaderShould("Content-Type", Equal("text/plain"))
			})
		})

		Context("when the target is /", func() {
			BeforeEach(func() {
				existingFile := path.Join(basePath, "one")
				Expect(createTextFile(existingFile, "1")).To(Succeed())
			})

			It("responds with 200 OK", func() {
				controller.Get(responseBuffer, "/")
				response = httptest.ParseResponse(responseBuffer)
				response.StatusShouldBe(200, "OK")
			})
		})
	})

	Describe("#Head", func() {
		var (
			getResponseBuffer = &bytes.Buffer{}
			getResponse       *httptest.ResponseMessage
		)

		BeforeEach(func() {
			controller.Get(getResponseBuffer, "/missing.txt")
			getResponse = httptest.ParseResponse(getResponseBuffer)

			controller.Head(responseBuffer, "/missing.txt")
			response = httptest.ParseResponse(responseBuffer)
		})

		It("responds with the same status as #Get would have", func() {
			getResponse.StatusShouldBe(404, "Not Found")
			response.StatusShouldBe(404, "Not Found")
		})
		It("responds with the same headers as #Get would have", func() {
			response.HeaderShould("Content-Type", Equal("text/plain"))
			response.HeaderShould("Content-Length", Equal("23"))
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
