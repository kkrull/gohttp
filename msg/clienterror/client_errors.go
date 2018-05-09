package clienterror

import (
	"fmt"
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

type NotFound struct {
	Path string
	body string
}

func (notFound *NotFound) Name() string {
	return "A resource that can not be found"
}

func (notFound *NotFound) WriteTo(client io.Writer) error {
	notFound.WriteHeader(client)
	msg.WriteBody(client, notFound.body)
	return nil
}

func (notFound *NotFound) WriteHeader(client io.Writer) error {
	msg.WriteStatus(client, NotFoundStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	notFound.body = fmt.Sprintf("Not found: %s", notFound.Path)
	msg.WriteContentLengthHeader(client, len(notFound.body))
	msg.WriteEndOfMessageHeader(client)
	return nil
}
