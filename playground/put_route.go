package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/success"
)

func NewPuttableRoute(path string) *PuttableRoute {
	return &PuttableRoute{
		Path:     path,
		Resource: &PuttableResource{},
	}
}

type PuttableRoute struct {
	Path     string
	Resource *PuttableResource
}

func (route *PuttableRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Resource)
}

// A resource where it's OK to PUT
type PuttableResource struct{}

func (*PuttableResource) Name() string {
	return "Puttable"
}

func (resource *PuttableResource) Post(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}

func (resource *PuttableResource) Put(client io.Writer, message http.RequestMessage) {
	success.RespondOkWithoutBody(client)
}
