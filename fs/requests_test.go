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
})
