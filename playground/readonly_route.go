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

func (route *ReadOnlyRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Target != "/method_options2" {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type ReadOnlyResource interface {
	Name() string
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
}

// Handles various read requests, but doesn't actually do anything
type ReadableNopResource struct{}

func (controller *ReadableNopResource) Name() string {
	return "Readonly NOP"
}

func (controller *ReadableNopResource) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *ReadableNopResource) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}
