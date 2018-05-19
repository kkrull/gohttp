package capability

import (
	"io"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

// Reports on server capabilities that are defined during startup and do not change after that
type StaticCapabilityServer struct {
	AvailableMethods []string
}

func (controller *StaticCapabilityServer) Options(client io.Writer) {
	msg.RespondWithAllowHeader(client, success.OKStatus, controller.AvailableMethods)
}
