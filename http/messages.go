package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

func NewGetMessage(path string) RequestMessage {
	return &requestMessage{
		method: "GET",
		target: path,
		path:   path,
	}
}

func NewHeadMessage(path string) RequestMessage {
	return &requestMessage{
		method: "HEAD",
		target: path,
		path:   path,
	}
}

func NewOptionsMessage(target string) RequestMessage {
	return &requestMessage{
		method: "OPTIONS",
		target: target,
		path:   target,
	}
}

func NewPutMessage(path string) RequestMessage {
	return &requestMessage{
		method: "PUT",
		target: path,
		path:   path,
	}
}

func NewTraceMessage(path string) RequestMessage {
	return &requestMessage{
		method: "TRACE",
		target: path,
		path:   path,
	}
}

func NewRequestMessage(method, path string) RequestMessage {
	return &requestMessage{
		method: method,
		target: path,
		path:   path,
	}
}

type requestMessage struct {
	method          string
	path            string
	target          string
	queryParameters []QueryParameter
}

func (message *requestMessage) Method() string {
	return message.method
}

func (message *requestMessage) Path() string {
	return message.path
}

func (message *requestMessage) AddQueryFlag(name string) {
	message.queryParameters = append(message.queryParameters, QueryParameter{Name: name})
}

func (message *requestMessage) AddQueryParameter(name, value string) {
	message.queryParameters = append(message.queryParameters, QueryParameter{Name: name, Value: value})
}

func (message *requestMessage) QueryParameters() []QueryParameter {
	return message.queryParameters
}

func (message *requestMessage) Target() string {
	return message.target
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
