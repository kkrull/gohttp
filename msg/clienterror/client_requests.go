package clienterror

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

func MethodNotAllowed(supportedMethods ...string) *methodNotAllowed {
	return &methodNotAllowed{SupportedMethods: supportedMethods}
}

type methodNotAllowed struct {
	SupportedMethods []string
}

func (notAllowed *methodNotAllowed) Handle(client io.Writer) error {
	msg.WriteStatusLine(client, 405, "Method Not Allowed")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", strings.Join(notAllowed.SupportedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
	return nil
}
