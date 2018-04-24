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
	var resources = map[string]bool{
		"/coffee": true,
		"/tea":    true,
	}

	_, ownResource := resources[requested.Target]
	if !ownResource {
		return nil
	}

	return http.MakeResourceRequest(requested, route.Resource)
}

type Teapot interface {
	Name() string
	Get(client io.Writer, target string)
}
