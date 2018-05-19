package clienterror

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

func MethodNotAllowed(supportedMethods ...string) *methodNotAllowed {
	return &methodNotAllowed{SupportedMethods: supportedMethods}
}

type methodNotAllowed struct {
	SupportedMethods []string
}

func (notAllowed *methodNotAllowed) Handle(client io.Writer) error {
	msg.RespondWithAllowHeader(client, MethodNotAllowedStatus, notAllowed.SupportedMethods)
	return nil
}
