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
		return route.routeToMethod(requested, route.Writable)
	case "/method_options2":
		return route.routeToMethod(requested, route.Readable)
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
	methods := map[string]Method{
		"GET":     &getMethod{},
		"HEAD":    &headMethod{},
		"OPTIONS": &optionsMethod{},
		"POST":    &postMethod{},
		"PUT":     &putMethod{},
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

type headRequest struct {
	Resource HeadResource
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Resource.Head(client)
	return nil
}

type HeadResource interface {
	Head(client io.Writer)
}

/* OPTIONS */

type optionsMethod struct{}

func (*optionsMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	supportedResource, ok := resource.(OptionsResource)
	if ok {
		return &optionsRequest{Resource: supportedResource}
	}

	return nil
}

type optionsRequest struct {
	Resource OptionsResource
}

func (request *optionsRequest) Handle(client io.Writer) error {
	request.Resource.Options(client)
	return nil
}

type OptionsResource interface {
	Options(client io.Writer)
}

/* POST */

type postMethod struct{}

func (*postMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	supportedResource, ok := resource.(PostResource)
	if ok {
		return &postRequest{Resource: supportedResource}
	}

	return nil
}

type postRequest struct {
	Resource PostResource
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Resource.Post(client)
	return nil
}

type PostResource interface {
	Post(client io.Writer)
}

/* PUT */

type putMethod struct{}

func (*putMethod) MakeRequest(requested *http.RequestLine, resource interface{}) http.Request {
	supportedResource, ok := resource.(PutResource)
	if ok {
		return &putRequest{Resource: supportedResource}
	}

	return nil
}

type putRequest struct {
	Resource PutResource
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Resource.Put(client)
	return nil
}

type PutResource interface {
	Put(client io.Writer)
}
