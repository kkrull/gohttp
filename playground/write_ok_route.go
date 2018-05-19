package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewWriteOKRoute(path string) *WriteOKRoute {
	return &WriteOKRoute{
		Path:     path,
		Resource: &WriteOKResource{},
	}
}

type WriteOKRoute struct {
	Path     string
	Resource *WriteOKResource
}

func (route *WriteOKRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

// A resource where it's OK to write
type WriteOKResource struct{}

func (*WriteOKResource) Name() string {
	return "Write OK"
}

func (resource *WriteOKResource) Post(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}

func (resource *WriteOKResource) Put(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
