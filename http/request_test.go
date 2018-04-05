package http_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	"bytes"
	"bufio"
	"github.com/kkrull/gohttp/http"
	"strings"
)

var _ = Describe("RFC7230RequestParser", func() {
	FDescribe("#ParseRequest", func() {
		var (
			parser  *http.RFC7230RequestParser
			request *http.Request
			err     *http.ParseError
		)

		BeforeEach(func() {
			parser = &http.RFC7230RequestParser{}
		})

		XIt("playground", func() {
			buffer := bytes.NewBufferString("abcd")
			bufioReader := bufio.NewReader(buffer)

			firstRound, firstReadErr := bufioReader.ReadString('b')
			for i, c := range firstRound {
				fmt.Printf("FIRST ROUND Byte %02d: %02x %2s\n", i, c, string(c))
			}
			if firstReadErr != nil {
				fmt.Printf("First round %s\n", firstReadErr.Error())
			}

			secondRound, secondReadErr := bufioReader.ReadString('e')
			for i, c := range secondRound {
				fmt.Printf("SECOND ROUND Byte %02d: %02x %2s\n", i, c, string(c))
			}
			if secondReadErr != nil {
				fmt.Printf("Second round %s\n", secondReadErr.Error())
			}
		})

		It("parses a well-formed request", func() {
			request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n\r\n"))
			Expect(request).To(BeEquivalentTo(&http.Request{
				Method:  "GET",
				Target:  "/",
				Version: "HTTP/1.1",
			}))
		})

		It("returns 400 Bad Request for an empty request", func() {
			buffer := bytes.NewBufferString("")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})

		It("returns 400 Bad Request for a request not containing 3 fields in the request line", func() {
			buffer := bytes.NewBufferString("\r\n")
			request, err = parser.ParseRequest(bufio.NewReader(buffer))
			Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
		})
	})

	Describe("#ParseRequest", func() {
		var (
			parser  http.RequestParser
			request *http.Request
			err     *http.ParseError
		)

		BeforeEach(func() {
			parser = http.RFC7230RequestParser{}
		})

		Describe("Section 3.1.1", func() {
			Context("given a well formed request line that is not too long", func() {
				It("parses the HTTP method as all text before the first space", func() {
					request, err = parser.ParseRequest(makeReader("GET / HTTP/1.1\r\n\r\n"))
					Expect(request).To(BeEquivalentTo(&http.Request{
						Method:  "GET",
						Target:  "/",
						Version: "HTTP/1.1",
					}))
				})
			})

			Context("when the request starts with whitespace", func() {
				It("Returns 400 Bad Request", func() {
					request, err = parser.ParseRequest(makeReader(" GET / HTTP/1.1\r\n\r\n"))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 400, Reason: "Bad Request"}))
				})
			})

			Context("when the request method is longer than 8,000 octets", func() {
				It("Returns 501 Not Implemented", func() {
					enormousMethod := strings.Repeat("POST", 2000) + "!"
					request, err = parser.ParseRequest(makeReader("%s / HTTP/1.1\r\n\r\n", enormousMethod))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 501, Reason: "Not Implemented"}))
				})
			})

			Context("when the request target is longer than 8,000 octets", func() {
				It("Returns 414 URI Too Long", func() {
					enormousTarget := strings.Repeat("/foo", 2000) + "/"
					request, err = parser.ParseRequest(makeReader("GET %s HTTP/1.1\r\n\r\n", enormousTarget))
					Expect(err).To(BeEquivalentTo(&http.ParseError{StatusCode: 414, Reason: "URI Too Long"}))
				})
			})
		})
	})
})

func makeReader(template string, values ...string) *bufio.Reader {
	text := fmt.Sprintf(template, values)
	return bufio.NewReader(bytes.NewBufferString(text))
}
