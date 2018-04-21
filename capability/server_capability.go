package capability

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

// Reports on server capabilities that are defined during startup and do not change after that
type StaticCapabilityController struct {
	AvailableMethods []string
}

func (controller *StaticCapabilityController) Options(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteHeader(client, "Allow", strings.Join(controller.AvailableMethods, ","))
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
