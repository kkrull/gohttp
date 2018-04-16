// HTTP 4xx Client Error responses from RFC 7231, Section 6.5
package clienterror

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kkrull/gohttp/msg"
)

type BadRequest struct {
	DisplayText string
}

func (badRequest BadRequest) WriteTo(client io.Writer) error {
	msg.WriteStatusLine(client, 400, "Bad Request")
	return nil
}

type NotFound struct {
	Target string
}

func (notFound NotFound) WriteTo(client io.Writer) error {
	msg.WriteStatusLine(client, 404, "Not Found")
	msg.WriteHeader(client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", notFound.Target)
	msg.WriteHeader(client, "Content-Length", strconv.Itoa(len(message)))
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, message)
	return nil
}
