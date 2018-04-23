package playground

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type StatelessOptionController struct {
}

func (controller *StatelessOptionController) Options(client io.Writer, target string) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", "GET,HEAD,POST,OPTIONS,PUT")
	msg.WriteEndOfMessageHeader(client)
}
