package log

import "github.com/kkrull/gohttp/http"

func NewLogRoute(path string) http.Route {
	return &Route{}
}

type Route struct{}

func (route *Route) Route(requested http.RequestMessage) http.Request {
	panic("implement me")
}
