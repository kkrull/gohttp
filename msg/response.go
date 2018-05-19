package msg

import (
	"io"
	"strings"
)

func RespondWithAllowHeader(client io.Writer, status Status, allowedMethods []string) {
	WriteStatus(client, status)
	WriteContentLengthHeader(client, 0)
	WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	WriteEndOfMessageHeader(client)
}
