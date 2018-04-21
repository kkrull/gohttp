package opt

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type StaticCapabilityController struct {
}

func (controller *StaticCapabilityController) Options(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
