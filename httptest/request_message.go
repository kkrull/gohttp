package httptest

import (
	"github.com/kkrull/gohttp/http"

	. "github.com/onsi/gomega"
)

type RequestMessage struct {
	MethodReturns      string
	PathReturns        string
	TargetReturns      string
	HeaderLinesReturns []string

	MakeResourceRequestReturns  http.Request
	makeResourceRequestReceived http.Resource

	queryParameters []http.QueryParameter
}

func (message *RequestMessage) Method() string {
	return message.MethodReturns
}

func (message *RequestMessage) Path() string {
	return message.PathReturns
}

func (message *RequestMessage) Target() string {
	return message.TargetReturns
}

func (message *RequestMessage) AddQueryParameter(name string, value string) {
	message.queryParameters = append(message.queryParameters, http.QueryParameter{Name: name, Value: value})
}

func (message *RequestMessage) QueryParameters() []http.QueryParameter {
	return message.queryParameters
}

func (message *RequestMessage) HeaderLines() []string {
	return message.HeaderLinesReturns
}

func (message *RequestMessage) MakeResourceRequest(resource http.Resource) http.Request {
	message.makeResourceRequestReceived = resource
	return message.MakeResourceRequestReturns
}

func (message *RequestMessage) MakeResourceRequestShouldHaveReceived(resource http.Resource) {
	ExpectWithOffset(1, message.makeResourceRequestReceived).To(BeIdenticalTo(resource))
}
