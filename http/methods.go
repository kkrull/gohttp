package http

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

/* GET */

type getMethod struct{}

func (method *getMethod) MakeRequest(requested *RequestLine, resource Resource) Request {
	supportedResource, ok := resource.(GetResource)
	if ok {
		return &getRequest{Resource: supportedResource}
	}

	return nil
}

type getRequest struct {
	Resource GetResource
}

func (request *getRequest) Handle(client io.Writer) error {
	request.Resource.Get(client)
	return nil
}

type GetResource interface {
	Get(client io.Writer)
}

/* HEAD */

type headMethod struct {
	Resource HeadResource
}

func (*headMethod) MakeRequest(requested *RequestLine, resource Resource) Request {
	supportedResource, ok := resource.(HeadResource)
	if ok {
		return &headRequest{Resource: supportedResource}
	}

	return nil
}

type headRequest struct {
	Resource HeadResource
}

func (request *headRequest) Handle(client io.Writer) error {
	request.Resource.Head(client)
	return nil
}

type HeadResource interface {
	Head(client io.Writer)
}

/* OPTIONS */

type optionsRequest struct {
	SupportedMethods []string
}

func (request *optionsRequest) Handle(client io.Writer) error {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", strings.Join(request.SupportedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
	return nil
}

/* POST */

type postMethod struct{}

func (*postMethod) MakeRequest(requested *RequestLine, resource Resource) Request {
	supportedResource, ok := resource.(PostResource)
	if ok {
		return &postRequest{Resource: supportedResource}
	}

	return nil
}

type postRequest struct {
	Resource PostResource
}

func (request *postRequest) Handle(client io.Writer) error {
	request.Resource.Post(client)
	return nil
}

type PostResource interface {
	Post(client io.Writer)
}

/* PUT */

type putMethod struct{}

func (*putMethod) MakeRequest(requested *RequestLine, resource Resource) Request {
	supportedResource, ok := resource.(PutResource)
	if ok {
		return &putRequest{Resource: supportedResource}
	}

	return nil
}

type putRequest struct {
	Resource PutResource
}

func (request *putRequest) Handle(client io.Writer) error {
	request.Resource.Put(client)
	return nil
}

type PutResource interface {
	Put(client io.Writer)
}
