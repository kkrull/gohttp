package playground

import (
	"io"
	"sort"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
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
	}

	return lookupRequest(requested, route.Resource)
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
	}

	return lookupRequest(requested, route.Resource)
}

type ReadWriteResource interface {
	Get(client io.Writer)
	Head(client io.Writer)
	Post(client io.Writer)
	Put(client io.Writer)
}

func lookupRequest(requested *http.RequestLine, resource interface{}) http.Request {
	if requested.Method == "OPTIONS" {
		return &knownOptionsRequest{
			SupportedMethods: supportedMethods(requested.Target, resource),
		}
	}

	method := knownMethods[requested.Method]
	if method == nil {
		return unknownHttpMethod(requested, resource)
	}

	request := method.MakeRequest(requested, resource)
	if request == nil {
		return unsupportedMethod(requested, resource)
	}

	return request
}

func unknownHttpMethod(requested *http.RequestLine, resource interface{}) http.Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func unsupportedMethod(requested *http.RequestLine, resource interface{}) http.Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func supportedMethods(target string, resource interface{}) []string {
	supported := []string{"OPTIONS"}
	for name, method := range knownMethods {
		imaginaryRequest := &http.RequestLine{Method: name, Target: target}
		request := method.MakeRequest(imaginaryRequest, resource)
		if request != nil {
			supported = append(supported, name)
		}
	}

	sort.Strings(supported)
	return supported
}

var knownMethods = map[string]Method{
	"GET":  &getMethod{},
	"HEAD": &headMethod{},
	"POST": &postMethod{},
	"PUT":  &putMethod{},
}

type Method interface {
	MakeRequest(requested *http.RequestLine, resource interface{}) http.Request
}
