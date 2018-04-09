package http_test

import (
	"github.com/kkrull/gohttp/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRequest", func() {
	It("accepts a base path", func() {
		request := http.GetRequest{
			BaseDirectory: "/tmp",
			Target:  "/",
			Version: "HTTP/1.1"}
		Expect(request).NotTo(BeNil())
	})
})
