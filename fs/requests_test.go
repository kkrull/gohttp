package fs_test

import (
	"bytes"
	"os"
	"path"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadOnlyFileSystem", func() {
	var (
		controller *fs.ReadOnlyFileSystem
		basePath   string

		response       *httptest.ResponseMessage
		responseBuffer *bytes.Buffer
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("ReadOnlyFileSystem", os.ModePerm)
		controller = &fs.ReadOnlyFileSystem{BaseDirectory: basePath}
		responseBuffer = &bytes.Buffer{}
	})

	Describe("#Get", func() {
		Context("when the resolved path does not exist", func() {
			BeforeEach(func() {
				controller.Get(responseBuffer, http.NewGetMessage("/missing.txt"))
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

		Describe("reading a file", func() {
			Context("when the path is a readable text file in the base path", func() {
				BeforeEach(func() {
					existingFile := path.Join(basePath, "readable.txt")
					Expect(createTextFile(existingFile, "A")).To(Succeed())
					controller.Get(responseBuffer, http.NewGetMessage("/readable.txt"))
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

			Context("when the path is a readable file named with a registered extension", func() {
				BeforeEach(func() {
					existingFile := path.Join(basePath, "image.jpeg")
					Expect(createTextFile(existingFile, "A")).To(Succeed())
				})

				It("sets Content-Type to the MIME type registered for that extension", func() {
					controller.Get(responseBuffer, http.NewGetMessage("/image.jpeg"))
					response = httptest.ParseResponse(responseBuffer)
					response.HeaderShould("Content-Type", Equal("image/jpeg"))
				})
			})

			Context("when the path is a readable file without an extension", func() {
				BeforeEach(func() {
					existingFile := path.Join(basePath, "assumed-text")
					Expect(createTextFile(existingFile, "A")).To(Succeed())
				})

				It("sets Content-Type to text/plain", func() {
					controller.Get(responseBuffer, http.NewGetMessage("/assumed-text"))
					response = httptest.ParseResponse(responseBuffer)
					response.HeaderShould("Content-Type", Equal("text/plain"))
				})
			})

			XContext("when the request contains a Range header for a range that exists in the requested file", func() {
				BeforeEach(func() {
					existingFile := path.Join(basePath, "readable.txt")
					Expect(createTextFile(existingFile, "ABC")).To(Succeed())
					requestMessage := &httptest.RequestMessage{
						MethodReturns: "GET",
						TargetReturns: "/readable.txt",
						PathReturns:   "/readable.txt",
					}
					requestMessage.AddHeader("Range", "0-1")

					controller.Get(responseBuffer, requestMessage)
					response = httptest.ParseResponse(responseBuffer)
				})

				It("responds with 200 OK", func() {
					response.StatusShouldBe(206, "Partial Content")
				})
				It("sets Content-Length to the number of bytes in the file", func() {
					response.HeaderShould("Content-Length", Equal("2"))
				})
				It("sets Content-Range to the resolved location in and total size of the file, in bytes", func() {
					response.HeaderShould("Content-Range", Equal("bytes 0-1/3"))
				})
				It("writes the contents of the file to the message body", func() {
					response.BodyShould(Equal("AB"))
				})
			})
		})

		Describe("reading a directory", func() {
			Context("when the path is /", func() {
				BeforeEach(func() {
					existingFile := path.Join(basePath, "one")
					Expect(createTextFile(existingFile, "1")).To(Succeed())
				})

				It("responds with 200 OK", func() {
					controller.Get(responseBuffer, http.NewGetMessage("/"))
					response = httptest.ParseResponse(responseBuffer)
					response.StatusShouldBe(200, "OK")
				})
			})
		})
	})

	Describe("#Head", func() {
		var (
			getResponseBuffer = &bytes.Buffer{}
			getResponse       *httptest.ResponseMessage
		)

		BeforeEach(func() {
			controller.Get(getResponseBuffer, http.NewGetMessage("/missing.txt"))
			getResponse = httptest.ParseResponse(getResponseBuffer)

			controller.Head(responseBuffer, http.NewHeadMessage("/missing.txt"))
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
