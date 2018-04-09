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
			basePath string
		)

		BeforeEach(func() {
			basePath = makeEmptyTestDirectory("GetRequest", os.ModePerm)
		})

		Context("when the resolved target does not exist", func() {
			It("Responds with 404 Not Found", func() {
				buffer := &bytes.Buffer{}
				writer := bufio.NewWriter(buffer)
				request := http.GetRequest{
					BaseDirectory: basePath,
					Target:        "/missing.txt",
					Version:       "HTTP/1.1"}
				Expect(request.Handle(writer)).To(Succeed())
				Expect(buffer.String()).To(Equal("HTTP/1.1 404 Not Found\r\n"))
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
