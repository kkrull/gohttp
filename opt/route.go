package opt

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *Route {
	controller := &StaticCapabilityController{}
	return &Route{Controller: controller}
}

type Route struct {
	Controller ServerCapabilityController
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	if requested.Method == "OPTIONS" && requested.Target == "*" {
		return &OptionsRequest{Controller: route.Controller}
	}

	return nil
}

type OptionsRequest struct {
	Controller ServerCapabilityController
}

func (request *OptionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client)
	return nil
}

type ServerCapabilityController interface {
	Options(writer io.Writer)
}
