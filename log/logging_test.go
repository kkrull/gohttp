package log_test

import (
	"bytes"

	"github.com/kkrull/gohttp/httptest"
	"github.com/kkrull/gohttp/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextLogger", func() {
	Describe("#Parsed", func() {
		var (
			logger *log.TextLogger
			output *bytes.Buffer
		)

		BeforeEach(func() {
			output = &bytes.Buffer{}
			logger = log.NewBufferedRequestLogger()

			requestMessage := &httptest.RequestMessage{
				MethodReturns:  "GET",
				TargetReturns:  "/foo",
				VersionReturns: "HTTP/1.1",
			}
			requestMessage.AddHeader("Content-Type", "text/plain")

			logger.Parsed(requestMessage)
			logger.WriteTo(output)
		})

		It("writes the request method and target", func() {
			Expect(output.String()).To(ContainSubstring("GET /foo HTTP/1.1"))
		})
		XIt("writes each header line", func() {
			Expect(output.String()).To(MatchRegexp("Content-Type: text/plain"))
		})
	})
})
