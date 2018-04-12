package fs_test

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/kkrull/gohttp/fs"
	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("DirectoryListing", func() {
	Describe("#WriteTo", func() {
		var (
			listing http.Response
			message *HttpMessage
			err     error
		)

		Context("always", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{Files: []string{}}
				listing.WriteTo(output)
				message = &HttpMessage{Text: output.String()}
			})

			It("returns no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("responds with 200 OK", func() {
				message.StatusShouldBe(200, "OK")
			})
			It("sets Content-Length to the size of the message", func() {
				contentLength, err := message.HeaderAsInt("Content-Length")
				Expect(err).NotTo(HaveOccurred())
				Expect(contentLength).To(BeNumerically(">", 0))
			})
			It("sets Content-Type to text/html", func() {
				message.ShouldHaveHeader("Content-Type", Equal("text/html"))
			})
		})

		Context("given no file names", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{Files: []string{}}
				listing.WriteTo(output)
				message = &HttpMessage{Text: output.String()}
			})

			It("has an empty list of file names", func() {
				message.BodyShould(MatchRegexp(".*<ul>\\s*<[/]ul>.*"))
			})
		})

		Context("given 1 or more file names", func() {
			BeforeEach(func() {
				output := &bytes.Buffer{}
				listing = &fs.DirectoryListing{
					Files:      []string{"one", "two"},
					HrefPrefix: "/files"}
				listing.WriteTo(output)
				message = &HttpMessage{Text: output.String()}
			})

			It("lists links to the files, using absolute paths with the given prefix", func() {
				message.BodyShould(ContainSubstring("<a href=\"/files/one\">one</a>"))
				message.BodyShould(ContainSubstring("<a href=\"/files/two\">two</a>"))
			})
		})
	})
})

type HttpMessage struct {
	Text string
}

func (message *HttpMessage) StatusShouldBe(status int, reason string) {
	Expect(message.Text).To(HavePrefix("HTTP/1.1 %d %s\r\n", status, reason))
}

func (message *HttpMessage) BodyShould(matcher types.GomegaMatcher) {
	_, body := message.splitMessageHeaderAndBody()
	Expect(body).To(matcher)
}

func (message HttpMessage) HeaderAsInt(name string) (int, error) {
	headers := message.headerFields()
	return strconv.Atoi(headers[name])
}

func (message *HttpMessage) ShouldHaveHeader(name string, match types.GomegaMatcher) {
	value := message.headerFields()[name]
	Expect(value).To(match)
}

func (message HttpMessage) headerFields() map[string]string {
	const indexAfterStartLine = 1
	messageHeader, _ := message.splitMessageHeaderAndBody()
	headerLines := strings.Split(messageHeader, "\r\n")[indexAfterStartLine:]
	headers := make(map[string]string)
	for _, line := range headerLines {
		field, value := parseHeader(line)
		headers[field] = value
	}

	return headers
}

func (message HttpMessage) splitMessageHeaderAndBody() (messageHeader, messageBody string) {
	const headerBodySeparator = "\r\n\r\n"
	split := strings.Split(message.Text, headerBodySeparator)
	return split[0], split[1]
}

func parseHeader(line string) (field, value string) {
	const optionalWhitespaceCharacters = " \t"
	fields := strings.Split(line, ":")
	field = fields[0]
	value = strings.Trim(fields[1], optionalWhitespaceCharacters)
	return
}
