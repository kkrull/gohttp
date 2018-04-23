package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *Route {
	return &Route{
		Controller: &StatelessOptionController{},
	}
}

type Route struct {
	Controller OptionController
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
	Controller OptionController
	Target     string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client, request.Target)
	return nil
}

type OptionController interface {
	Options(client io.Writer, target string)
}
