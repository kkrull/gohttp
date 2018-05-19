package msg

import (
	"bytes"
	"io"
)

func RespondInPlainText(client io.Writer, status Status, bodyText []byte) {
	WriteStatus(client, status)
	WriteContentTypeHeader(client, "text/plain")
	WriteContentLengthHeader(client, len(bodyText))
	WriteEndOfMessageHeader(client)
	CopyToBody(client, bytes.NewReader(bodyText))
}
