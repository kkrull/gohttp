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
		)

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				listing = &DirectoryListing{Files: []string{"one", "two"}}
				output = &bytes.Buffer{}
				listing.WriteTo(output)
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				Expect(output.String()).To(haveStatus(200, "OK"))
			})
			It("sets Content-Length to the size of the message", func() {
				parser := HttpMessageParser{Text: output.String()}
				contentLength, err := parser.HeaderAsInt("Content-Length")
				Expect(err).NotTo(HaveOccurred())
				Expect(contentLength).To(BeNumerically(">", 0))
			})
			It("sets Content-Type to text/html", func() {
				Expect(output.String()).To(containHeader("Content-Type", "text/html"))
			})
			It("lists links to the files in the base path", func() {
				Expect(output.String()).To(haveMessageBody("one\ntwo\n"))
			})
		})
	})
})

type HttpMessageParser struct {
	Text string
}

func (parser HttpMessageParser) HeaderAsInt(name string) (int, error) {
	headers := parser.headerFields()
	return strconv.Atoi(headers[name])
}

func (parser HttpMessageParser) headerFields() map[string]string {
	const indexAfterStartLine = 1
	headerLines := strings.Split(parser.messageHeader(), "\r\n")[indexAfterStartLine:]
	headers := make(map[string]string)
	for _, line := range headerLines {
		field, value := parseHeader(line)
		headers[field] = value
	}

	return headers
}

func (parser HttpMessageParser) messageHeader() string {
	const headerBodySeparator = "\r\n\r\n"
	return strings.Split(parser.Text, headerBodySeparator)[0]
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

func haveMessageBody(message string) types.GomegaMatcher {
	return HaveSuffix(fmt.Sprintf("\r\n\r\n%s", message))
}
