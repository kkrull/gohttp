package http

import (
	"fmt"
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

const (
	CONNECT string = "CONNECT"
	GET     string = "GET"
	HEAD    string = "HEAD"
	OPTIONS string = "OPTIONS"
	POST    string = "POST"
	PUT     string = "PUT"
	TRACE   string = "TRACE"
)

func NewGetMessage(path string) RequestMessage {
	return &requestMessage{
		method: GET,
		target: path,
		path:   path,
	}
}

func NewHeadMessage(path string) RequestMessage {
	return &requestMessage{
		method: HEAD,
		target: path,
		path:   path,
	}
}

// Creates an OPTIONS request to the specified target, which can either be a path starting with /
// or an asterisk-form query of the server as a whole (https://tools.ietf.org/html/rfc7230#section-5.3.4).
func NewOptionsMessage(targetAsteriskOrPath string) RequestMessage {
	return &requestMessage{
		method: OPTIONS,
		target: targetAsteriskOrPath,
		path:   targetAsteriskOrPath,
	}
}

func NewPutMessage(path string) RequestMessage {
	return &requestMessage{
		method: PUT,
		target: path,
		path:   path,
	}
}

func NewTraceMessage(path string) RequestMessage {
	return &requestMessage{
		method: TRACE,
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
	headers         []header
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

func (message *requestMessage) HeaderLines() []string {
	lines := make([]string, len(message.headers))
	for i, header := range message.headers {
		lines[i] = header.String()
	}

	return lines
}

func (message *requestMessage) HeaderValues(field string) []string {
	values := make([]string, 0)
	for _, header := range message.headers {
		if header.Field != field {
			continue
		}

		values = append(values, header.Value)
	}

	return values
}

func (message *requestMessage) AddHeader(field, value string) {
	message.headers = append(message.headers, header{Field: field, Value: value})
}

func (message *requestMessage) NotImplemented() Response {
	return &servererror.NotImplemented{Method: message.method}
}

func (message *requestMessage) MakeResourceRequest(resource Resource) Request {
	if message.method == OPTIONS {
		return &optionsRequest{
			SupportedMethods: message.supportedMethods(resource),
		}
	}

	method := knownMethods[message.method]
	if method == nil {
		return message.unknownHttpMethod(resource)
	}

	request, isSupported := method.MakeRequest(message, resource)
	if !isSupported {
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
	supported := []string{OPTIONS}
	for name, method := range knownMethods {
		imaginaryRequest := &requestMessage{method: name, target: message.target}
		_, isSupported := method.MakeRequest(imaginaryRequest, resource)
		if isSupported {
			supported = append(supported, name)
		}
	}

	sort.Strings(supported)
	return supported
}

var knownMethods = map[string]Method{
	GET:  &getMethod{},
	HEAD: &headMethod{},
	POST: &postMethod{},
	PUT:  &putMethod{},
}

type Method interface {
	MakeRequest(requested *requestMessage, resource Resource) (request Request, isSupported bool)
}

type header struct {
	Field string
	Value string
}

func (header *header) String() string {
	return fmt.Sprintf("%s: %s", header.Field, header.Value)
}

// Handles requests of supported HTTP methods for a resource
type Resource interface {
	Name() string
}
