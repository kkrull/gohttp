// Testing related to the http package
package httptest

import (
	"bytes"
	"strconv"
	"strings"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func ParseResponse(text *bytes.Buffer) *ResponseMessage {
	return &ResponseMessage{Text: text.String()}
}

type ResponseMessage struct {
	Text string
}

func (message *ResponseMessage) ShouldBeWellFormed() {
	Expect(message.Text).To(ContainSubstring("\r\n\r\n"))
}

func (message *ResponseMessage) StatusShouldBe(status int, reason string) {
	Expect(message.Text).To(HavePrefix("HTTP/1.1 %d %s\r\n", status, reason))
}

func (message *ResponseMessage) BodyShould(matcher types.GomegaMatcher) {
	_, body := message.splitMessageHeaderAndBody()
	Expect(body).To(matcher)
}

func (message ResponseMessage) HeaderAsInt(name string) (int, error) {
	headers := message.headerFields()
	return strconv.Atoi(headers[name])
}

func (message *ResponseMessage) HeaderShould(name string, match types.GomegaMatcher) {
	value := message.headerFields()[name]
	Expect(value).To(match)
}

func (message ResponseMessage) headerFields() map[string]string {
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

func (message ResponseMessage) splitMessageHeaderAndBody() (messageHeader, messageBody string) {
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
