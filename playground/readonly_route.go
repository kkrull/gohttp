package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewReadOnlyRoute(path string) *ReadOnlyRoute {
	return &ReadOnlyRoute{
		Path:     path,
		Resource: &ReadableNopResource{},
	}
}

type ReadOnlyRoute struct {
	Path     string
	Resource ReadOnlyResource
}

func (route *ReadOnlyRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadOnlyResource interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
	Head(client io.Writer, message http.RequestMessage)
}

// Handles various read requests, but doesn't actually do anything
type ReadableNopResource struct{}

func (controller *ReadableNopResource) Name() string {
	return "Readonly NOP"
}

func (controller *ReadableNopResource) Get(client io.Writer, message http.RequestMessage) {
	controller.Head(client, message)
}

func (controller *ReadableNopResource) Head(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
