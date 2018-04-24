package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() *Route {
	return &Route{
		Readable: &ReadableNopResource{},
		Writable: &ReadWriteNopResource{},
	}
}

type Route struct {
	Readable ReadOnlyResource
	Writable ReadWriteResource
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	switch requested.Target {
	case "/method_options":
		routedRequest := route.routeToMethod(requested, route.Writable)
		if routedRequest != nil {
			return routedRequest
		}

		switch requested.Method {
		case "HEAD", "OPTIONS", "POST", "PUT":
			return route.makeRequest(requested, route.Writable)
		default:
			return nil
		}
	case "/method_options2":
		routedRequest := route.routeToMethod(requested, route.Readable)
		if routedRequest != nil {
			return routedRequest
		}

		switch requested.Method {
		case "HEAD", "OPTIONS":
			return route.makeRequest(requested, route.Readable)
		default:
			return nil
		}
	default:
		return nil
	}
}

func (route *Route) routeToMethod(requested *http.RequestLine, resource interface{}) http.Request {
	if requested.Method != "GET" {
		return nil
	}

	getMethod := &getMethod{}
	request := getMethod.MakeRequest(requested, resource)
	if request != nil {
		return request
	}

	return nil
}

type getMethod struct {}

func (method *getMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	gettableResource, ok := resource.(GetResource)
	if ok {
		return &getRequest{
			Controller: gettableResource,
			Target:     requested.Target,
		}
	}

	return nil
}


func (route *Route) makeRequest(requested *http.RequestLine, controller interface{}) http.Request {
	//TODO KDK: Boil this down to a typecast for the desired controller/method
	switch requested.Method {
	case "GET":
		readController, _ := controller.(ReadOnlyResource)
		return &getRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "HEAD":
		readController, _ := controller.(ReadOnlyResource)
		return &headRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "OPTIONS":
		readController, _ := controller.(ReadOnlyResource)
		return &optionsRequest{
			Controller: readController,
			Target:     requested.Target,
		}
	case "POST":
		writeController, _ := controller.(ReadWriteResource)
		return &postRequest{
			Controller: writeController,
			Target:     requested.Target,
		}
	case "PUT":
		writeController, _ := controller.(ReadWriteResource)
		return &putRequest{
			Controller: writeController,
			Target:     requested.Target,
		}
	default:
		return nil
	}
}

type getRequest struct {
	Controller GetResource
	Target     string
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

type headRequest struct {
	Controller ReadOnlyResource
	Target     string
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}

type optionsRequest struct {
	Controller ReadOnlyResource
	Target     string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client, request.Target)
	return nil
}

type postRequest struct {
	Controller ReadWriteResource
	Target     string
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Controller.Post(client, request.Target)
	return nil
}

type putRequest struct {
	Controller ReadWriteResource
	Target     string
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Controller.Put(client, request.Target)
	return nil
}

type GetResource interface {
	Get(client io.Writer, target string)
}

type ReadOnlyResource interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Options(client io.Writer, target string)
}

type ReadWriteResource interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Options(client io.Writer, target string)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}
