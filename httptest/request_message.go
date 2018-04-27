package httptest

import "github.com/kkrull/gohttp/http"

type RequestMessage struct {
	ItsMethod       string
	ItsPath         string
	ItsTarget       string
	queryParameters []http.QueryParameter
}

func (*RequestMessage) MakeResourceRequest(resource http.Resource) http.Request {
	panic("implement me")
}

func (message *RequestMessage) Method() string {
	return message.ItsMethod
}

func (message *RequestMessage) Target() string {
	return message.ItsTarget
}

func (message *RequestMessage) AddQueryParameter(name string, value string) {
	message.queryParameters = append(message.queryParameters, http.QueryParameter{Name: name, Value: value})
}

func (message *RequestMessage) QueryParameters() []http.QueryParameter {
	return message.queryParameters
}
