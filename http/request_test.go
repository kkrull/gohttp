package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"fmt"
	"github.com/kkrull/gohttp/http"
	"bytes"
	"bufio"
)

var _ = Describe("RFC7230RequestParser", func() {
	var parser http.RequestParser
	BeforeEach(func() {
		parser = http.RFC7230RequestParser{}
	})

	Context("when the request starts with whitespace", func() {
		It("Responds 400 Bad Request", func() {
			_, err := parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})
	})

	Context("when the request method is longer than 8,000 octets", func() {
		It("responds 501 Not Implemented", func() {
			enormousMethod := strings.Repeat("POST", 2000) + "!"
			_, err := parser.ParseRequest(makeReader(fmt.Sprintf("%s / HTTP/1.1\r\n\r\n", enormousMethod)))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 501, Reason: "Not Implemented"}))
		})
	})

	Context("when the request target is longer than 8,000 octets", func() {
		It("responds 414 URI Too Long", func() {
			enormousTarget := strings.Repeat("/foo", 2000) + "/"
			_, err := parser.ParseRequest(makeReader(fmt.Sprintf("GET %s HTTP/1.1\r\n\r\n", enormousTarget)))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 414, Reason: "URI Too Long"}))
		})
	})

	XIt("Section 3 paragraph 3 (encoding must be a superset of US-ASCII)")
	XIt("Section 3 paragraph 5 (recipient must reject request with whitespace between the start-line and the first header)")
})

func makeReader(text string) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(text))
}
