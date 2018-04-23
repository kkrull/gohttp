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
	case "/method_options":
		switch requested.Method {
		case "GET":
			return &getRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "HEAD":
			return &headRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "OPTIONS":
			return &optionsRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "POST":
			return &postRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "PUT":
			return &putRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		default:
			return nil
		}
	case "/method_options2":
		switch requested.Method {
		case "GET":
			return &getRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "HEAD":
			return &headRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		case "OPTIONS":
			return &optionsRequest{
				Controller: route.Controller,
				Target:     requested.Target,
			}
		default:
			return nil
		}
	default:
		return nil
	}
}

type getRequest struct {
	Controller Controller
	Target     string
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

type headRequest struct {
	Controller Controller
	Target     string
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}

type optionsRequest struct {
	Controller Controller
	Target     string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client, request.Target)
	return nil
}

type postRequest struct {
	Controller Controller
	Target     string
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Controller.Post(client, request.Target)
	return nil
}

type putRequest struct {
	Controller Controller
	Target     string
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Controller.Put(client, request.Target)
	return nil
}

type Controller interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Options(client io.Writer, target string)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}
