package clienterror

import (
	"fmt"
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

func RespondMethodNotAllowed(client io.Writer, allowedMethods []string) {
	msg.WriteStatus(client, MethodNotAllowedStatus)
	msg.WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
}

func RespondNotFound(client io.Writer, path string) {
	msg.WriteStatus(client, NotFoundStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("Not found: %s", path)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, body)
}
