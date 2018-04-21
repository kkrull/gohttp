package teapot

import "github.com/kkrull/gohttp/http"

func NewRoute() http.Route {
	controller := &IdentityController{}
	return &Route{Controller: controller}
}

type Route struct {
	Controller Controller
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	if requested.Target == "/coffee" {
		return &GetRequest{
			Controller: route.Controller,
			Target:     requested.Target,
		}
	}

	return nil
}
