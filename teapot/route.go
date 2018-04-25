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

func (route *Route) Route(requested *http.RequestLine) http.Request {
	if !route.Teapot.RespondsTo(requested.Target) {
		return nil
	}

	return http.MakeResourceRequest(requested, route.Teapot)
}

type Teapot interface {
	Name() string
	Get(client io.Writer, target string)
	RespondsTo(target string) bool
}
