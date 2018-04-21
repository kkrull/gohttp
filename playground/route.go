package playground

import "github.com/kkrull/gohttp/http"

func NewRoute() *Route {
	return &Route{}
}

type Route struct {
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	return nil
}
