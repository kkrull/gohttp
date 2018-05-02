// HTTP 4xx Client Error responses from RFC 7231, Section 6.5
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
	msg.WriteStatusLine(client, 400, "Bad Request")
	//msg.WriteEndOfMessageHeader(client)
	return nil
}

type NotFound struct {
	Path string
	body string
}

func (notFound *NotFound) WriteTo(client io.Writer) error {
	notFound.WriteHeader(client)
	msg.WriteBody(client, notFound.body)
	return nil
}

func (notFound *NotFound) WriteHeader(client io.Writer) error {
	msg.WriteStatusLine(client, 404, "Not Found")
	msg.WriteContentTypeHeader(client, "text/plain")

	notFound.body = fmt.Sprintf("Not found: %s", notFound.Path)
	msg.WriteContentLengthHeader(client, len(notFound.body))
	msg.WriteEndOfMessageHeader(client)
	return nil
}
