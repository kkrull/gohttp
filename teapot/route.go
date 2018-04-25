package teapot

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewRoute() http.Route {
	teapot := &IdentityTeapot{}
	return &Route{Resource: teapot}
}

type Route struct {
	Resource Teapot
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	if !route.Resource.RespondsTo(requested.Target) {
		return nil
	}

	return http.MakeResourceRequest(requested, route.Resource)
}

type Teapot interface {
	Name() string
	Get(client io.Writer, target string)
	RespondsTo(target string) bool
}
