package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *Route {
	return &Route{
		Controller: &AllowedMethodsController{},
	}
}

type Route struct {
	Controller Controller
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	switch requested.Target {
	case "/method_options", "/method_options2":
		return &optionsRequest{
			Controller: route.Controller,
			Target:     requested.Target,
		}
	default:
		return nil
	}
}

type optionsRequest struct {
	Controller Controller
	Target     string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client, request.Target)
	return nil
}

type Controller interface {
	Options(client io.Writer, target string)
}
