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

func (route *Route) Route(message http.RequestMessage) http.Request {
	if !route.Teapot.RespondsTo(message.Path()) {
		return nil
	}

	return message.MakeResourceRequest(route.Teapot)
}

type Teapot interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
	RespondsTo(path string) bool
}
