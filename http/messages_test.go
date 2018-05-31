package http_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("requestMessage", func() {
	Describe("#MakeResourceRequest", func() {
		Context("given a message with the PATCH method", func() {
			var (
				message  http.RequestMessage
				resource *ResourceMock
				request  http.Request
			)

			BeforeEach(func() {
				message = http.NewRequestMessage("PATCH", "/existing")
				resource = &ResourceMock{}
				request = message.MakeResourceRequest(resource)
			})

			It("returns a request that calls PatchResource#Patch", func() {
				request.Handle(&bytes.Buffer{})
				resource.PatchShouldHaveBeenCalled("/existing")
			})
		})
	})
})
