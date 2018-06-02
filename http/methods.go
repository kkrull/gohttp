package http

import (
	"io"

	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

/* DELETE */

type deleteMethod struct{}

func (method *deleteMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(DeleteResource)
	if ok {
		return &deleteRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type deleteRequest struct {
	Message  RequestMessage
	Resource DeleteResource
}

func (request *deleteRequest) Handle(client io.Writer) error {
	request.Resource.Delete(client, request.Message)
	return nil
}

type DeleteResource interface {
	Delete(client io.Writer, message RequestMessage)
}

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

type optionsMethod struct{}

func (*optionsMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(OptionsResource)
	if ok {
		return &dynamicOptionsRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

// Asks an OptionsResource what it thinks its supported HTTP methods are,
// when there is no single answer that is known ahead of time
type dynamicOptionsRequest struct {
	Message  RequestMessage
	Resource OptionsResource
}

func (request *dynamicOptionsRequest) Handle(client io.Writer) error {
	request.Resource.Options(client, request.Message)
	return nil
}

type OptionsResource interface {
	Options(client io.Writer, message RequestMessage)
}

// Responds with a static set of supported HTTP methods that are known a priori
type staticOptionsRequest struct {
	SupportedMethods []string
}

func (request *staticOptionsRequest) Handle(client io.Writer) error {
	msg.RespondWithAllowHeader(client, success.OKStatus, request.SupportedMethods)
	return nil
}

/* PATCH */

type patchMethod struct{}

func (*patchMethod) MakeRequest(message *requestMessage, resource Resource) (request Request, isSupported bool) {
	supportedResource, ok := resource.(PatchResource)
	if ok {
		return &patchRequest{
			Message:  message,
			Resource: supportedResource,
		}, true
	}

	return nil, false
}

type patchRequest struct {
	Message  RequestMessage
	Resource PatchResource
}

func (request *patchRequest) Handle(client io.Writer) error {
	request.Resource.Patch(client, request.Message)
	return nil
}

type PatchResource interface {
	Patch(client io.Writer, message RequestMessage)
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
