package clienterror

import (
	"fmt"
	"io"

	"github.com/kkrull/gohttp/msg"
)

func RespondMethodNotAllowed(client io.Writer, allowedMethods []string) {
	msg.RespondWithAllowHeader(client, MethodNotAllowedStatus, allowedMethods)
}

func RespondNotFound(client io.Writer, path string) {
	msg.WriteStatus(client, NotFoundStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("Not found: %s", path)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, body)
}
