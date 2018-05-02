package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRedirectRoute() http.Route {
	return &RedirectRoute{
		Resource: &GoBackHomeResource{},
	}
}

type RedirectRoute struct {
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

func (*GoBackHomeResource) Get(client io.Writer, req http.RequestMessage) {
	panic("implement me")
}
