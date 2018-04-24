package http

import (
	"sort"

	"github.com/kkrull/gohttp/msg/clienterror"
)

func MakeResourceRequest(requested *RequestLine, resource interface{}) Request {
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

func unknownHttpMethod(requested *RequestLine, resource interface{}) Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func unsupportedMethod(requested *RequestLine, resource interface{}) Request {
	return clienterror.MethodNotAllowed(supportedMethods(requested.Target, resource)...)
}

func supportedMethods(target string, resource interface{}) []string {
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
	MakeRequest(requested *RequestLine, resource interface{}) Request
}
