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
	} else if requested.Method == "OPTIONS" {
		return &knownOptionsRequest{
			SupportedMethods: supportedMethods(requested.Target, route.Resource),
		}
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
		return &knownOptionsRequest{
			SupportedMethods: supportedMethods(requested.Target, route.Resource),
		}
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
	method := knownMethods[requested.Method]
	if method == nil {
		return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
	}

	request := method.MakeRequest(requested, resource)
	if request == nil {
		return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
	}

	return request
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
