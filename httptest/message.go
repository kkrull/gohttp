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
	ExpectWithOffset(1, message.Text).To(HavePrefix("HTTP/1.1"), "missing status-line")
	ExpectWithOffset(1, message.Text).To(
		ContainSubstring("\r\n\r\n"),
		"missing CRLF{2} to mark the end of the message header")
}

func (message *ResponseMessage) StatusShouldBe(status int, reason string) {
	ExpectWithOffset(1, message.Text).To(HavePrefix("HTTP/1.1 %d %s\r\n", status, reason))
}

func (message *ResponseMessage) BodyShould(matcher types.GomegaMatcher) {
	_, body := message.splitMessageHeaderAndBody()
	ExpectWithOffset(1, body).To(matcher)
}

func (message ResponseMessage) HeaderAsInt(name string) (int, error) {
	headers := message.headerFields()
	return strconv.Atoi(headers[name])
}

func (message *ResponseMessage) HeaderShould(name string, match types.GomegaMatcher) {
	value, ok := message.headerFields()[name]
	ExpectWithOffset(1, ok).To(BeTrue(), "%s does not exist", name)
	ExpectWithOffset(1, value).To(match, name)
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
