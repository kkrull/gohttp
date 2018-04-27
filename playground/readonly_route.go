package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewReadOnlyRoute() *ReadOnlyRoute {
	return &ReadOnlyRoute{Resource: &ReadableNopResource{}}
}

type ReadOnlyRoute struct {
	Resource ReadOnlyResource
}

func (route *ReadOnlyRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/method_options2" {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadOnlyResource interface {
	Name() string
	Get(client io.Writer, req http.RequestMessage)
	Head(client io.Writer, target string)
}

// Handles various read requests, but doesn't actually do anything
type ReadableNopResource struct{}

func (controller *ReadableNopResource) Name() string {
	return "Readonly NOP"
}

func (controller *ReadableNopResource) Get(client io.Writer, req http.RequestMessage) {
	controller.Head(client, req.Path())
}

func (controller *ReadableNopResource) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}
