package opt

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *ServerCapabilityRoute {
	controller := &StaticCapabilityController{
		AvailableMethods: []string{"GET", "HEAD"},
	}

	return &ServerCapabilityRoute{Controller: controller}
}

type ServerCapabilityRoute struct {
	Controller ServerCapabilityController
}

func (route *ServerCapabilityRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Method == "OPTIONS" && requested.Target == "*" {
		return &optionsRequest{Controller: route.Controller}
	}

	return nil
}

type optionsRequest struct {
	Controller ServerCapabilityController
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client)
	return nil
}

// Reports the global, generic capabilities of this server, without regard to resource or state
type ServerCapabilityController interface {
	Options(writer io.Writer)
}
