package success

import (
	"bytes"
	"io"

	"github.com/kkrull/gohttp/msg"
)

func RespondOKWithKnownBody(client io.Writer, contentType string, body []byte) {
	msg.WriteStatus(client, OKStatus)
	msg.WriteContentTypeHeader(client, contentType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.CopyToBody(client, bytes.NewReader(body))
}

func RespondOkWithoutBody(client io.Writer) {
	msg.WriteStatus(client, OKStatus)
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
