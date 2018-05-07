package capability

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

// Reports on server capabilities that are defined during startup and do not change after that
type StaticCapabilityServer struct {
	AvailableMethods []string
}

func (controller *StaticCapabilityServer) Options(client io.Writer) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteHeader(client, "Allow", strings.Join(controller.AvailableMethods, ","))
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
