package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
)

func NewReadWriteRoute() *ReadWriteRoute {
	return &ReadWriteRoute{
		Resource: &ReadWriteNopResource{},
	}
}

type ReadWriteRoute struct {
	Resource ReadWriteResource
}

func (route *ReadWriteRoute) Route(requested *http.RequestLine) http.Request {
	if requested.TheTarget != "/method_options" {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadWriteResource interface {
	Name() string
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}

// Handles various read/write requests, but doesn't actually do anything
type ReadWriteNopResource struct{}

func (controller *ReadWriteNopResource) Name() string {
	return "Read/Write NOP"
}

func (controller *ReadWriteNopResource) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *ReadWriteNopResource) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Post(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Put(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
