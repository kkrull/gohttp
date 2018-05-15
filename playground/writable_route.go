package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewReadWriteRoute() *ReadWriteRoute {
	return &ReadWriteRoute{
		Resource: &ReadWriteNopResource{},
	}
}

type ReadWriteRoute struct {
	Resource ReadWriteResource
}

func (route *ReadWriteRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/method_options" {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadWriteResource interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
	Head(client io.Writer, message http.RequestMessage)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}

// Handles various read/write requests, but doesn't actually do anything
type ReadWriteNopResource struct{}

func (controller *ReadWriteNopResource) Name() string {
	return "Read/Write NOP"
}

func (controller *ReadWriteNopResource) Get(client io.Writer, message http.RequestMessage) {
	controller.Head(client, message)
}

func (controller *ReadWriteNopResource) Head(client io.Writer, message http.RequestMessage) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Post(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Put(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
