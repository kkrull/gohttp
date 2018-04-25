package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
)

func MakeResourceRequest(requested *RequestLine, resource Resource) Request {
	if requested.Method == "OPTIONS" {
		return &optionsRequest{
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

func unknownHttpMethod(requested *RequestLine, resource Resource) Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func unsupportedMethod(requested *RequestLine, resource Resource) Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func supportedMethods(target string, resource Resource) []string {
	supported := []string{"OPTIONS"}
	for name, method := range knownMethods {
		imaginaryRequest := &RequestLine{Method: name, Target: target}
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
