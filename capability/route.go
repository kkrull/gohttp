package capability

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

func NewRoute() *ServerCapabilityRoute {
	controller := &StaticCapabilityServer{
		AvailableMethods: []string{http.GET, http.HEAD},
	}

	return &ServerCapabilityRoute{Controller: controller}
}

type ServerCapabilityRoute struct {
	Controller ServerResource
}

func (route *ServerCapabilityRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Target() != "*" {
		return nil
	} else if requested.Method() != http.OPTIONS {
		return clienterror.MethodNotAllowed(http.OPTIONS)
	}

	return &optionsRequest{Resource: route.Controller}
}

type optionsRequest struct {
	Resource ServerResource
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Resource.Options(client)
	return nil
}

// Reports the global, generic capabilities of this server, without regard to resource or state
type ServerResource interface {
	Options(writer io.Writer)
}
