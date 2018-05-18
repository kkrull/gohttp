package httptest

import (
	"fmt"

	"github.com/kkrull/gohttp/http"

	. "github.com/onsi/gomega"
)

type RequestMessage struct {
	MethodReturns string
	PathReturns   string
	TargetReturns string

	MakeResourceRequestReturns  http.Request
	makeResourceRequestReceived http.Resource

	queryParameters []http.QueryParameter
	headers         []header
	body            []byte
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

func (message *RequestMessage) AddHeader(field string, value string) {
	message.headers = append(message.headers, header{Name: field, Value: value})
}

func (message *RequestMessage) HeaderLines() []string {
	lines := make([]string, len(message.headers))
	for i, header := range message.headers {
		lines[i] = fmt.Sprintf("%s: %s", header.Name, header.Value)
	}

	return lines
}

func (message *RequestMessage) HeaderValues(field string) (values []string) {
	values = make([]string, 0)
	for _, header := range message.headers {
		values = append(values, header.Value)
	}

	return values
}

func (message *RequestMessage) Body() []byte {
	return message.body
}

func (message *RequestMessage) MakeResourceRequest(resource http.Resource) http.Request {
	message.makeResourceRequestReceived = resource
	return message.MakeResourceRequestReturns
}

func (message *RequestMessage) MakeResourceRequestShouldHaveReceived(resource http.Resource) {
	ExpectWithOffset(1, message.makeResourceRequestReceived).To(BeIdenticalTo(resource))
}

type header struct {
	Name, Value string
}
