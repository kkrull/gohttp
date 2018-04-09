package http_test

import (
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

		It("accepts a base path", func() {
			request := http.GetRequest{
				BaseDirectory: "/tmp",
				Target:        "/",
				Version:       "HTTP/1.1"}
			Expect(request).NotTo(BeNil())
		})
	})
})

func makeEmptyTestDirectory(testName string, fileMode os.FileMode) string {
	testPath := path.Join(".test", testName)
	Expect(os.RemoveAll(testPath)).To(Succeed())
	Expect(os.MkdirAll(testPath, fileMode)).To(Succeed())
	return testPath
}
