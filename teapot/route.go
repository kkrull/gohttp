package teapot

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

func NewRoute() http.Route {
	controller := &IdentityController{}
	return &Route{Controller: controller}
}

type Route struct {
	Controller Controller
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	var resources = map[string]http.Request{
		"/coffee": &GetCoffeeRequest{Controller: route.Controller},
		"/tea":    &GetTeaRequest{Controller: route.Controller},
	}

	request, ownResource := resources[requested.Target]
	if !ownResource {
		return nil
	} else if requested.Method != "GET" {
		return clienterror.MethodNotAllowed("GET")
	}

	return request
}

type Controller interface {
	GetCoffee(client io.Writer)
	GetTea(client io.Writer)
}
