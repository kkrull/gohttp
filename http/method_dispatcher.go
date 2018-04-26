package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

type RequestLine struct {
	Method          string
	Target          string
	QueryParameters map[string]string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}

func (requestLine *RequestLine) MakeResourceRequest(resource Resource) Request {
	if requestLine.Method == "OPTIONS" {
		return &optionsRequest{
			SupportedMethods: requestLine.supportedMethods(resource),
		}
	}

	method := knownMethods[requestLine.Method]
	if method == nil {
		return requestLine.unknownHttpMethod(resource)
	}

	request := method.MakeRequest(requestLine, resource)
	if request == nil {
		return requestLine.unsupportedMethod(resource)
	}

	return request
}

func (requestLine *RequestLine) unknownHttpMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(requestLine.supportedMethods(resource)...)
}

func (requestLine *RequestLine) unsupportedMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(requestLine.supportedMethods(resource)...)
}

func (requestLine *RequestLine) supportedMethods(resource Resource) []string {
	supported := []string{"OPTIONS"}
	for name, method := range knownMethods {
		imaginaryRequest := &RequestLine{Method: name, Target: requestLine.Target}
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
	MakeRequest(requested *RequestLine, resource Resource) Request
}

// Handles requests of supported HTTP methods for a resource
type Resource interface {
	Name() string
}
