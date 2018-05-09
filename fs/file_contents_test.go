package fs_test

import (
	"bytes"
	"os"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotFound", func() {
	var (
		resource       *clienterror.NotFound
		basePath       string
		response       *httptest.ResponseMessage
		responseBuffer = &bytes.Buffer{}
	)

	BeforeEach(func() {
		responseBuffer.Reset()
		basePath = makeEmptyTestDirectory("NotFound", os.ModePerm)
	})

	XDescribe("#Get", func() {
		Context("when the source path does not exist", func() {
			BeforeEach(func() {
				resource = &clienterror.NotFound{Path: "/missing.txt"}
				//resource.Get(responseBuffer, http.NewGetMessage("/missing.txt"))
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
	})
})
