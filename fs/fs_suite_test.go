package fs_test

import (
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestFs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fs Suite")
}

/* HttpMessage */

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

func (message *HttpMessage) HeaderShould(name string, match types.GomegaMatcher) {
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
