package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

func NewRequestMessage(method, target string) RequestMessage {
	return &RequestLine{
		TheMethod: method,
		TheTarget: target,
	}
}

type RequestLine struct {
	TheMethod          string
	TheTarget          string
	TheQueryParameters map[string]string
}

func (requestLine *RequestLine) Method() string {
	return requestLine.TheMethod
}

func (requestLine *RequestLine) Target() string {
	return requestLine.TheTarget
}

func (requestLine *RequestLine) QueryParameters() map[string]string {
	return requestLine.TheQueryParameters
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.TheMethod}
}

func (requestLine *RequestLine) MakeResourceRequest(resource Resource) Request {
	if requestLine.TheMethod == "OPTIONS" {
		return &optionsRequest{
			SupportedMethods: requestLine.supportedMethods(resource),
		}
	}

	method := knownMethods[requestLine.TheMethod]
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
		imaginaryRequest := &RequestLine{TheMethod: name, TheTarget: requestLine.TheTarget}
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
