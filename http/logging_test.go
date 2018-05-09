package http_test

import (
	"bytes"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TextLogger", func() {
	Describe("#Parsed", func() {
		var (
			logger http.RequestLogger
			output *bytes.Buffer
		)

		BeforeEach(func() {
			output = &bytes.Buffer{}
			logger = http.TextLogger{Writer: output}

			requestMessage := &httptest.RequestMessage{
				MethodReturns:      "GET",
				TargetReturns:      "/foo",
				HeaderLinesReturns: []string{"Content-Type: text/plain"},
			}

			logger.Parsed(requestMessage)
		})

		It("writes the request method and target", func() {
			Expect(output.String()).To(ContainSubstring("GET /foo"))
		})
		It("writes each header line", func() {
			Expect(output.String()).To(MatchRegexp("Content-Type: text/plain"))
		})
	})
})
