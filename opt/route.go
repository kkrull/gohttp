package opt

import "github.com/kkrull/gohttp/http"

func NewRoute() http.Route {
	return &Route{}
}

type Route struct {
	Controller Controller
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	if requested.Method == "OPTIONS" && requested.Target == "*" {
		return &OptionsRequest{Controller: route.Controller}
	}

	return nil
}
