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
	if requested.Target != "/method_options2" {
		return nil
	} else if requested.Method == "OPTIONS" {
		return optionsRequest(requested, route.Resource)
	}

	return resourceRequest(requested, route.Resource)
}

type ReadOnlyResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
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
	if requested.Target != "/method_options" {
		return nil
	} else if requested.Method == "OPTIONS" {
		return optionsRequest(requested, route.Resource)
	}

	return resourceRequest(requested, route.Resource)
}

type ReadWriteResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Post(client io.Writer)
	Put(client io.Writer)
}

func resourceRequest(requested *http.RequestLine, resource interface{}) http.Request {
	methods := map[string]Method{
		"GET":  &getMethod{},
		"HEAD": &headMethod{},
		"POST": &postMethod{},
		"PUT":  &putMethod{},
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

func optionsRequest(requested *http.RequestLine, resource interface{}) http.Request {
	methods := map[string]Method{
		"GET":  &getMethod{},
		"HEAD": &headMethod{},
		"POST": &postMethod{},
		"PUT":  &putMethod{},
	}

	supportedMethods := []string{"OPTIONS"}
	for name, method := range methods {
		imaginaryRequest := &http.RequestLine{Method: name, Target: requested.Target}
		request := method.MakeRequest(imaginaryRequest, resource)
		if request != nil {
			supportedMethods = append(supportedMethods, name)
		}
	}

	return &knownOptionsRequest{SupportedMethods: supportedMethods}
}

type Method interface {
	MakeRequest(requested *http.RequestLine, resource interface{}) http.Request
}
