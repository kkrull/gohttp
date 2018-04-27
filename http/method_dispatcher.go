package http

import (
	"sort"
	"strings"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

func NewGetMessage(target string) RequestMessage {
	return &requestMessage{
		method: "GET",
		target: target,
	}
}

func NewHeadMessage(target string) RequestMessage {
	return &requestMessage{
		method: "HEAD",
		target: target,
	}
}

func NewOptionsMessage(target string) RequestMessage {
	return &requestMessage{
		method: "OPTIONS",
		target: target,
	}
}

func NewPutMessage(target string) RequestMessage {
	return &requestMessage{
		method: "PUT",
		target: target,
	}
}

func NewTraceMessage(target string) RequestMessage {
	return &requestMessage{
		method: "TRACE",
		target: target,
	}
}

func NewRequestMessage(method, target string) RequestMessage {
	return &requestMessage{
		method: method,
		target: target,
	}
}

type requestMessage struct {
	method string
	target string
}

func (message *requestMessage) Method() string {
	return message.method
}

func (message *requestMessage) Path() string {
	fields := strings.Split(message.target, "?")
	return fields[0]
}

func (message *requestMessage) Target() string {
	return message.target
}

func (message *requestMessage) QueryParameters() []QueryParameter {
	panic("implement me")
}

func (message *requestMessage) NotImplemented() Response {
	return &servererror.NotImplemented{Method: message.method}
}

func (message *requestMessage) MakeResourceRequest(resource Resource) Request {
	if message.method == "OPTIONS" {
		return &optionsRequest{
			SupportedMethods: message.supportedMethods(resource),
		}
	}

	method := knownMethods[message.method]
	if method == nil {
		return message.unknownHttpMethod(resource)
	}

	request := method.MakeRequest(message, resource)
	if request == nil {
		return message.unsupportedMethod(resource)
	}

	return request
}

func (message *requestMessage) unknownHttpMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(message.supportedMethods(resource)...)
}

func (message *requestMessage) unsupportedMethod(resource Resource) Request {
	return clienterror.MethodNotAllowed(message.supportedMethods(resource)...)
}

func (message *requestMessage) supportedMethods(resource Resource) []string {
	supported := []string{"OPTIONS"}
	for name, method := range knownMethods {
		imaginaryRequest := &requestMessage{method: name, target: message.target}
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
