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
		case "OPTIONS", "POST", "PUT":
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
		case "OPTIONS":
			return route.makeRequest(requested, route.Readable)
		default:
			return nil
		}
	default:
		return nil
	}
}

type ReadOnlyResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Options(client io.Writer)
}

type ReadWriteResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Options(client io.Writer)
	Post(client io.Writer)
	Put(client io.Writer)
}

func (route *Route) routeToMethod(requested *http.RequestLine, resource interface{}) http.Request {
	methods := map[string]Method {
		"GET": &getMethod{},
		"HEAD": &headMethod{},
	}

	method := methods[requested.Method]
	if method == nil {
		return nil
	}

	request := method.MakeRequest(requested, resource)
	if request != nil {
		return request
	}

	return nil
}

type Method interface {
	MakeRequest(requested *http.RequestLine, resource interface{}) http.Request
}

/* GET */

type getMethod struct{}

func (method *getMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	supportedResource, ok := resource.(GetResource)
	if ok {
		return &getRequest{Resource: supportedResource}
	}

	return nil
}

type getRequest struct {
	Resource GetResource
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Resource.Get(client)
	return nil
}

type GetResource interface {
	Get(client io.Writer)
}

/* HEAD */

type headMethod struct {
	Resource HeadResource
}

func (*headMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	supportedResource, ok := resource.(HeadResource)
	if ok {
		return &headRequest{Resource: supportedResource}
	}

	return nil
}

type HeadResource interface {
	Head(client io.Writer)
}

type headRequest struct {
	Resource HeadResource
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Resource.Head(client)
	return nil
}

/* OPTIONS */

type optionsRequest struct {
	Controller ReadOnlyResource
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client)
	return nil
}

/* POST */

type postRequest struct {
	Controller ReadWriteResource
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Controller.Post(client)
	return nil
}

/* PUT */

type putRequest struct {
	Controller ReadWriteResource
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Controller.Put(client)
	return nil
}









func (route *Route) makeRequest(requested *http.RequestLine, controller interface{}) http.Request {
	//TODO KDK: Boil this down to a typecast for the desired controller/method
	switch requested.Method {
	case "OPTIONS":
		readController, _ := controller.(ReadOnlyResource)
		return &optionsRequest{Controller: readController}
	case "POST":
		writeController, _ := controller.(ReadWriteResource)
		return &postRequest{Controller: writeController}
	case "PUT":
		writeController, _ := controller.(ReadWriteResource)
		return &putRequest{Controller: writeController}
	default:
		return nil
	}
}
