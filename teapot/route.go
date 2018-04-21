package teapot

import "github.com/kkrull/gohttp/http"

func NewRoute() http.Route {
	return &teapotRoute{}
}

type teapotRoute struct {
}

func (route *teapotRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Target == "/coffee" {
		return &GetCoffeeRequest{}
	}

	return nil
}
