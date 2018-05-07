package teapot

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() http.Route {
	teapot := &IdentityTeapot{}
	return &Route{Teapot: teapot}
}

type Route struct {
	Teapot Teapot
}

func (route *Route) Route(requested http.RequestMessage) http.Request {
	if !route.Teapot.RespondsTo(requested.Path()) {
		return nil
	}

	return requested.MakeResourceRequest(route.Teapot)
}

type Teapot interface {
	Name() string
	Get(client io.Writer, req http.RequestMessage)
	RespondsTo(path string) bool
}
