package clienterror_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/msg/clienterror"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("MethodNotAllowed", func() {
	var (
		request  http.Request
		response = &bytes.Buffer{}
	)

	Describe("#Handle", func() {
		BeforeEach(func() {
			response.Reset()
			request = &clienterror.MethodNotAllowed{SupportedMethods: []string{"GET", "HEAD"}}
			request.Handle(response)
		})

		It("responds 405 Method Not Allowed", httptest.ShouldHaveNoBody(response, 405, "Method Not Allowed"))
		It("sets Allow to GET and HEAD", httptest.ShouldAllowMethods(response, "GET", "HEAD"))
	})
})
