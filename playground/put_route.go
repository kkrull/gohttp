package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewNopPutRoute(path string) *NopPutRoute {
	return &NopPutRoute{
		Path:     path,
		Resource: &NopPutResource{},
	}
}

type NopPutRoute struct {
	Path     string
	Resource *NopPutResource
}

func (route *NopPutRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

// A resource where it's OK to PUT
type NopPutResource struct{}

func (*NopPutResource) Name() string {
	return "NOP Put"
}

func (resource *NopPutResource) Put(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
