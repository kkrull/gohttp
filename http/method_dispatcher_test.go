package http_test

import (
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestMessage", func() {
	var message http.RequestMessage

	Describe("#Path", func() {
		It("returns the full target, for a target without a query or fragment", func() {
			message = http.NewGetMessage("/widget")
			Expect(message.Path()).To(Equal("/widget"))
		})
		It("returns the part of the target before any ?", func() {
			message = http.NewGetMessage("/widget?field=value")
			Expect(message.Path()).To(Equal("/widget"))
		})
	})
})
