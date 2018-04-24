package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewReadOnlyRoute() *ReadOnlyRoute {
	return &ReadOnlyRoute{Resource: &ReadableNopResource{}}
}

type ReadOnlyRoute struct {
	Resource ReadOnlyResource
}

func (route *ReadOnlyRoute) Route(requested *http.RequestLine) http.Request {
	switch requested.Target {
	case "/method_options2":
		return routeToMethod(requested, route.Resource)
	default:
		return nil
	}
}

type ReadOnlyResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Options(client io.Writer)
}

func NewReadWriteRoute() *ReadWriteRoute {
	return &ReadWriteRoute{
		Resource: &ReadWriteNopResource{},
	}
}

type ReadWriteRoute struct {
	Resource ReadWriteResource
}

func (route *ReadWriteRoute) Route(requested *http.RequestLine) http.Request {
	switch requested.Target {
	case "/method_options":
		return routeToMethod(requested, route.Resource)
	default:
		return nil
	}
}

type ReadWriteResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Options(client io.Writer)
	Post(client io.Writer)
	Put(client io.Writer)
}

func routeToMethod(requested *http.RequestLine, resource interface{}) http.Request {
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
