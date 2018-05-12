package http

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

/* GET */

type getMethod struct{}

func (method *getMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(GetResource)
	if ok {
		return &getRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type getRequest struct {
	Message  RequestMessage
	Resource GetResource
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Resource.Get(client, request.Message)
	return nil
}

type GetResource interface {
	Get(client io.Writer, message RequestMessage)
}

/* HEAD */

type headMethod struct{}

func (*headMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(HeadResource)
	if ok {
		return &headRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type headRequest struct {
	Message  RequestMessage
	Resource HeadResource
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Resource.Head(client, request.Message)
	return nil
}

type HeadResource interface {
	Head(client io.Writer, message RequestMessage)
}

/* OPTIONS */

type optionsRequest struct {
	SupportedMethods []string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", strings.Join(request.SupportedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
	return nil
}

/* POST */

type postMethod struct{}

func (*postMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(PostResource)
	if ok {
		return &postRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type postRequest struct {
	Message  RequestMessage
	Resource PostResource
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Resource.Post(client, request.Message)
	return nil
}

type PostResource interface {
	Post(client io.Writer, message RequestMessage)
}

/* PUT */

type putMethod struct{}

func (*putMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(PutResource)
	if ok {
		return &putRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type putRequest struct {
	Message  RequestMessage
	Resource PutResource
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Resource.Put(client, request.Message)
	return nil
}

type PutResource interface {
	Put(client io.Writer, message RequestMessage)
}
