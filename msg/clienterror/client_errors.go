package clienterror

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type BadRequest struct {
	DisplayText string
}

func (badRequest *BadRequest) WriteTo(client io.Writer) error {
	return badRequest.WriteHeader(client)
}

func (badRequest *BadRequest) WriteHeader(client io.Writer) error {
	msg.WriteStatus(client, BadRequestStatus)
	//msg.WriteEndOfMessageHeader(client)
	return nil
}
