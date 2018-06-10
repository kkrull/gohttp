package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/redirect"
)

func NewRedirectRoute(path string) *RedirectRoute {
	return &RedirectRoute{
		Path:     path,
		Resource: &GoBackHomeResource{},
	}
}

type RedirectRoute struct {
	Path     string
	Resource RelocatedResource
}

func (route *RedirectRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/redirect" {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

type RelocatedResource interface {
	Name() string
}

type GoBackHomeResource struct{}

func (*GoBackHomeResource) Name() string {
	return "Relocated Resource"
}

func (*GoBackHomeResource) Get(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, redirect.FoundStatus)
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Location", "/")
	msg.WriteEndOfMessageHeader(client)
}
