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

	return http.MakeResourceRequest(requested, route.Resource)
}

type ReadOnlyResource interface {
	Name() string
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
}

func NewReadWriteRoute() *ReadWriteRoute {
	return &ReadWriteRoute{
		Resource: &ReadWriteNopResource{},
	}
}

type ReadWriteRoute struct {
	Resource ReadWriteResource
}

func (route *ReadWriteRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Target != "/method_options" {
		return nil
	}

	return http.MakeResourceRequest(requested, route.Resource)
}

type ReadWriteResource interface {
	Name() string
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
	Post(client io.Writer, target string)
	Put(client io.Writer, target string)
}
