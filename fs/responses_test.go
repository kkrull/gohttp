package fs_test

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	. "github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("DirectoryListing", func() {
	Describe("#WriteTo", func() {
		var (
			listing http.Response
			output  *bytes.Buffer
			err     error
			parser  *HttpMessageParser
		)

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				listing = &DirectoryListing{Files: []string{"one", "two"}}
				output = &bytes.Buffer{}
				listing.WriteTo(output)
				parser = &HttpMessageParser{Text: output.String()}
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(output.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the size of the message", func() {
				contentLength, err := parser.HeaderAsInt("Content-Length")
				Expect(err).NotTo(HaveOccurred())
				Expect(contentLength).To(BeNumerically(">", 0))
			})
			It("sets Content-Type to text/html", func() {
				Expect(output.String()).To(containHeader("Content-Type", "text/html"))
			})
			It("lists links to the files in the base path", func() {
				parser.BodyShouldContain("<a href=\"/one\">one</a>")
				parser.BodyShouldContain("<a href=\"/two\">two</a>")
			})
		})
	})
})

type HttpMessageParser struct {
	Text string
}

func (parser HttpMessageParser) BodyShouldContain(substring string) {
	_, body := parser.splitMessageHeaderAndBody()
	Expect(body).To(ContainSubstring(substring))
}

func (parser HttpMessageParser) HeaderAsInt(name string) (int, error) {
	headers := parser.headerFields()
	return strconv.Atoi(headers[name])
}

func (parser HttpMessageParser) headerFields() map[string]string {
	const indexAfterStartLine = 1
	messageHeader, _ := parser.splitMessageHeaderAndBody()
	headerLines := strings.Split(messageHeader, "\r\n")[indexAfterStartLine:]
	headers := make(map[string]string)
	for _, line := range headerLines {
		field, value := parseHeader(line)
		headers[field] = value
	}

	return headers
}

func (parser HttpMessageParser) splitMessageHeaderAndBody() (messageHeader, messageBody string) {
	const headerBodySeparator = "\r\n\r\n"
	split := strings.Split(parser.Text, headerBodySeparator)
	return split[0], split[1]
}

func parseHeader(line string) (field, value string) {
	const optionalWhitespaceCharacters = " \t"
	fields := strings.Split(line, ":")
	field = fields[0]
	value = strings.Trim(fields[1], optionalWhitespaceCharacters)
	return
}

func haveStatus(status int, reason string) types.GomegaMatcher {
	return HavePrefix(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, reason))
}

func containHeader(name string, value string) types.GomegaMatcher {
	return ContainSubstring(fmt.Sprintf("%s: %s\r\n", name, value))
}
