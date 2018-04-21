package opt

import "github.com/kkrull/gohttp/http"

func NewRoute() *Route {
	controller := &StaticCapabilitiesController{}
	return &Route{Controller: controller}
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
