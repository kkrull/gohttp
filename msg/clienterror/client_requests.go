package clienterror

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

type MethodNotAllowed struct {
	SupportedMethods []string
}

func (notAllowed *MethodNotAllowed) Handle(client io.Writer) error {
	msg.WriteStatusLine(client, 405, "Method Not Allowed")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", strings.Join(notAllowed.SupportedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
	return nil
}
