package http_test

import (
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("requestMessage", func() {
	Describe("#MakeResourceRequest", func() {
		Context("given a message with a PATCH method", func() {
			XIt("returns a request that calls PatchResource#Patch", func() {
				message := http.NewRequestMessage("PATCH", "/existing")
				resource := &PatchResourceMock{}
				request := message.MakeResourceRequest(resource)
				_, ok := request.(http.PatchResource)
				Expect(ok).To(BeTrue(), "Request of type %T to be compatible", request)
			})
		})
	})
})
