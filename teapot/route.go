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
	if requested.Method != "GET" {
		return clienterror.MethodNotAllowed("GET")
	}

	switch requested.Target {
	case "/coffee":
		return &GetCoffeeRequest{
			Controller: route.Controller,
		}
	case "/tea":
		return &GetTeaRequest{
			Controller: route.Controller,
		}
	default:
		return nil
	}
}

type Controller interface {
	GetCoffee(client io.Writer)
	GetTea(client io.Writer)
}
