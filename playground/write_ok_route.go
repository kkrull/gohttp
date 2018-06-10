package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewNopPostRoute(path string) *NopPostRoute {
	return &NopPostRoute{
		Path:     path,
		Resource: &NopPostResource{},
	}
}

type NopPostRoute struct {
	Path     string
	Resource *NopPostResource
}

func (route *NopPostRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

// A resource where it's OK to write
type NopPostResource struct{}

func (*NopPostResource) Name() string {
	return "NOP Post"
}

func (resource *NopPostResource) Post(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
