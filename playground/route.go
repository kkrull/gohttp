package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *Route {
	return &Route{
		ReadController:  &ReadableNopController{},
		WriteController: &WritableNopController{},
	}
}

type Route struct {
	ReadController  ReadableController
	WriteController ReadWriteController
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	switch requested.Target {
	case "/method_options":
		switch requested.Method {
		case "GET", "HEAD", "OPTIONS", "POST", "PUT":
			return route.makeRequest(requested, route.WriteController)
		default:
			return nil
		}
	case "/method_options2":
		switch requested.Method {
		case "GET", "HEAD", "OPTIONS":
			return route.makeRequest(requested, route.ReadController)
		default:
			return nil
		}
	default:
		return nil
	}
}

func (route *Route) makeRequest(requested *http.RequestLine, controller interface{}) http.Request {
	//TODO KDK: Boil this down to a typecast for the desired controller/method
	switch requested.Method {
	case "GET":
		readController, _ := controller.(ReadableController)
		return &getRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "HEAD":
		readController, _ := controller.(ReadableController)
		return &headRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "OPTIONS":
		readController, _ := controller.(ReadableController)
		return &optionsRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "POST":
		writeController, _ := controller.(ReadWriteController)
		return &postRequest{
			Controller: writeController,
			Target:     requested.Target,
		}
	case "PUT":
		writeController, _ := controller.(ReadWriteController)
		return &putRequest{
			Controller: writeController,
			Target:     requested.Target,
		}
	default:
		return nil
	}
}

type getRequest struct {
	Controller ReadableController
	Target     string
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

type headRequest struct {
	Controller ReadableController
	Target     string
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}

type optionsRequest struct {
	Controller ReadableController
	Target     string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client, request.Target)
	return nil
}

type postRequest struct {
	Controller ReadWriteController
	Target     string
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Controller.Post(client, request.Target)
	return nil
}

type putRequest struct {
	Controller ReadWriteController
	Target     string
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Controller.Put(client, request.Target)
	return nil
}

type ReadableController interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Options(client io.Writer, target string)
}

type ReadWriteController interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Options(client io.Writer, target string)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}
