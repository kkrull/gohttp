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

var _ = Describe("ExistingFile", func() {
	var (
		resource     *fs.ExistingFile
		basePath     string
		existingFile string

		response       *httptest.ResponseMessage
		responseBuffer = &bytes.Buffer{}
	)

	BeforeEach(func() {
		basePath = makeEmptyTestDirectory("ExistingFile", os.ModePerm)
		responseBuffer.Reset()
	})

	Describe("#Get", func() {
		Context("when the path is a readable file in the base path, and there are no Range headers", func() {
			BeforeEach(func() {
				existingFile = path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "A")).To(Succeed())

				resource = &fs.ExistingFile{Filename: existingFile}
				resource.Get(responseBuffer, http.NewGetMessage("/readable.txt"))
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
				existingFile = path.Join(basePath, "image.jpeg")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to the MIME type registered for that extension", func() {
				resource = &fs.ExistingFile{Filename: existingFile}
				resource.Get(responseBuffer, http.NewGetMessage("/image.jpeg"))
				response = httptest.ParseResponse(responseBuffer)
				response.HeaderShould("Content-Type", Equal("image/jpeg"))
			})
		})

		Context("when the path is a readable file without an extension", func() {
			BeforeEach(func() {
				existingFile = path.Join(basePath, "assumed-text")
				Expect(createTextFile(existingFile, "A")).To(Succeed())
			})

			It("sets Content-Type to text/plain", func() {
				resource = &fs.ExistingFile{Filename: existingFile}
				resource.Get(responseBuffer, http.NewGetMessage("/assumed-text"))
				response = httptest.ParseResponse(responseBuffer)
				response.HeaderShould("Content-Type", Equal("text/plain"))
			})
		})

		Context("when the request contains a single Range header with single byte range that exists in the requested file", func() {
			BeforeEach(func() {
				existingFile = path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "ABC")).To(Succeed())
				requestMessage := &httptest.RequestMessage{
					MethodReturns: "GET",
					TargetReturns: "/readable.txt",
					PathReturns:   "/readable.txt",
				}
				requestMessage.AddHeader("Range", "bytes=0-1")

				resource = &fs.ExistingFile{Filename: existingFile}
				resource.Get(responseBuffer, requestMessage)
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds with a well-formed message", func() {
				response.ShouldBeWellFormed()
			})
			It("responds with 200 OK", func() {
				response.StatusShouldBe(206, "Partial Content")
			})
			It("sets Content-Type to the MIME type registered for that extension", func() {
				response.HeaderShould("Content-Type", HavePrefix("text/plain"))
			})
			It("sets Content-Length to the number of bytes in the selected range(s)", func() {
				response.HeaderShould("Content-Length", Equal("2"))
			})
			It("sets Content-Range to the resolved location in and total size of the file, in bytes", func() {
				response.HeaderShould("Content-Range", Equal("bytes 0-1/3"))
			})
			It("writes the contents of the file to the message body", func() {
				response.BodyShould(Equal("AB"))
			})
		})

		Context("when the request contains 2 or more Range headers", func() {
			BeforeEach(func() {
				existingFile = path.Join(basePath, "readable.txt")
				Expect(createTextFile(existingFile, "ABCD")).To(Succeed())

				requestMessage := &httptest.RequestMessage{
					MethodReturns: "GET",
					TargetReturns: "/readable.txt",
					PathReturns:   "/readable.txt",
				}
				requestMessage.AddHeader("Range", "bytes=0-1")
				requestMessage.AddHeader("Range", "bytes=2-3")

				resource = &fs.ExistingFile{Filename: existingFile}
				resource.Get(responseBuffer, requestMessage)
				response = httptest.ParseResponse(responseBuffer)
			})

			It("responds as if it were a request for the whole file – that's way too complicated to handle right now", func() {
				response.StatusShouldBe(200, "OK")
				response.HeaderShould("Content-Length", Equal("4"))
			})
		})
	})

	Describe("#Head", func() {
		BeforeEach(func() {
			existingFile = path.Join(basePath, "readable.txt")
			Expect(createTextFile(existingFile, "A")).To(Succeed())

			resource = &fs.ExistingFile{Filename: existingFile}
			resource.Head(responseBuffer, http.NewHeadMessage("/readable.txt"))
			response = httptest.ParseResponse(responseBuffer)
		})

		It("returns the same status as #Get", func() {
			response.StatusShouldBe(200, "OK")
		})
		It("sets the same headers as #Get", func() {
			response.HeaderShould("Content-Length", Equal("1"))
			response.HeaderShould("Content-Type", HavePrefix("text/plain"))
		})
		It("has no body", func() {
			response.BodyShould(BeEmpty())
		})
	})

	Describe("#Options", func() {
		BeforeEach(func() {
			resource = &fs.ExistingFile{Filename: "/bonafide.txt"}
			requestMessage := http.NewOptionsMessage("/bonafide.txt")
			request := requestMessage.MakeResourceRequest(resource)

			request.Handle(responseBuffer)
		})

		It("supports read operations", httptest.ShouldAllowMethods(responseBuffer, "GET", "HEAD", "OPTIONS"))
	})
})
