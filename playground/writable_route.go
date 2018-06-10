package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewReadWriteRoute(path string) *ReadWriteRoute {
	return &ReadWriteRoute{
		Path:     path,
		Resource: &ReadWriteNopResource{},
	}
}

type ReadWriteRoute struct {
	Path     string
	Resource ReadWriteResource
}

func (route *ReadWriteRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadWriteResource interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
	Head(client io.Writer, message http.RequestMessage)
	Post(client io.Writer, message http.RequestMessage)
	Put(client io.Writer, message http.RequestMessage)
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
	success.RespondOkWithoutBody(client)
}

func (controller *ReadWriteNopResource) Post(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}

func (controller *ReadWriteNopResource) Put(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
