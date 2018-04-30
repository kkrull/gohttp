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
	path, _, _ := message.splitTarget()
	return path
}

func (message *requestMessage) QueryParameters() []QueryParameter {
	_, query, _ := message.splitTarget()
	if query == "" {
		return nil
	}

	return parseQueryString(query)
}

func (message *requestMessage) splitTarget() (path, query, fragment string) {
	//there is a path

	splitOnQuery := strings.Split(message.target, "?")
	if len(splitOnQuery) == 1 {
		//there is a path
		//there is no query string
		//there may be a fragment
		splitOnFragment := strings.Split(splitOnQuery[0], "#")
		if len(splitOnFragment) == 1 {
			//there is a path
			//there is no query string
			//there is no fragment
			return message.target, "", ""
		}

		//there is a path
		//there is no query string
		//there is a fragment
		return splitOnFragment[0], "", splitOnFragment[1]
	}

	//there is a path
	//there is a query string
	//there may be a fragment
	splitOnFragment := strings.Split(splitOnQuery[1], "#")
	if len(splitOnFragment) == 1 {
		//there is a path
		//there is a query string
		//there is no fragment
		return splitOnQuery[0], splitOnQuery[1], ""
	}

	//there is a path
	//there is a query string
	//there is a fragment
	return splitOnQuery[0], splitOnFragment[0], splitOnFragment[1]
}

func parseQueryString(rawQuery string) []QueryParameter {
	stringParameters := strings.Split(rawQuery, "&")
	parsedParameters := make([]QueryParameter, len(stringParameters))
	for i, stringParameter := range stringParameters {
		parsedParameters[i] = parseQueryParameter(stringParameter)
	}

	return parsedParameters
}

func parseQueryParameter(stringParameter string) QueryParameter {
	nameValueFields := strings.Split(stringParameter, "=")
	if len(nameValueFields) == 1 {
		return QueryParameter{Name: nameValueFields[0]}
	} else {
		return QueryParameter{Name: nameValueFields[0], Value: nameValueFields[1]}
	}
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
