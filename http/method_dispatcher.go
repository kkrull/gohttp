package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

func NewGetMessage(target string) RequestMessage {
	return &requestMessage{
		TheMethod: "GET",
		TheTarget: target,
	}
}

func NewHeadMessage(target string) RequestMessage {
	return &requestMessage{
		TheMethod: "HEAD",
		TheTarget: target,
	}
}

func NewOptionsMessage(target string) RequestMessage {
	return &requestMessage{
		TheMethod: "OPTIONS",
		TheTarget: target,
	}
}

func NewPutMessage(target string) RequestMessage {
	return &requestMessage{
		TheMethod: "PUT",
		TheTarget: target,
	}
}

func NewTraceMessage(target string) RequestMessage {
	return &requestMessage{
		TheMethod: "TRACE",
		TheTarget: target,
	}
}

func NewRequestMessage(method, target string) RequestMessage {
	return &requestMessage{
		TheMethod: method,
		TheTarget: target,
	}
}

type requestMessage struct {
	TheMethod          string
	TheTarget          string
}

func (requestLine *requestMessage) Method() string {
	return requestLine.TheMethod
}

func (requestLine *requestMessage) Target() string {
	return requestLine.TheTarget
}

func (requestLine *requestMessage) QueryParameters() []QueryParameter {
	panic("implement me")
}

func (requestLine *requestMessage) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.TheMethod}
}

func (requestLine *requestMessage) MakeResourceRequest(resource Resource) Request {
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

func (requestLine *requestMessage) unknownHttpMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(requestLine.supportedMethods(resource)...)
}

func (requestLine *requestMessage) unsupportedMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(requestLine.supportedMethods(resource)...)
}

func (requestLine *requestMessage) supportedMethods(resource Resource) []string {
	supported := []string{"OPTIONS"}
	for name, method := range knownMethods {
		imaginaryRequest := &requestMessage{TheMethod: name, TheTarget: requestLine.TheTarget}
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
	MakeRequest(requested *requestMessage, resource Resource) Request
}

// Handles requests of supported HTTP methods for a resource
type Resource interface {
	Name() string
}
