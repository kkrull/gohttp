package teapot

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

func NewRoute() http.Route {
	controller := &IdentityTeapot{}
	return &Route{Resource: controller}
}

type Route struct {
	Resource Teapot
}

func (route *Route) Route(requested *http.RequestLine) http.Request {
	var resources = map[string]http.Request{
		"/coffee": &GetCoffeeRequest{Controller: route.Resource},
		"/tea":    &GetTeaRequest{Controller: route.Resource},
	}

	request, ownResource := resources[requested.Target]
	if !ownResource {
		return nil
	} else if requested.Method != "GET" {
		return clienterror.MethodNotAllowed("GET")
	}

	return request
}

type Teapot interface {
	Name() string
	Get(client io.Writer, target string)
	GetCoffee(client io.Writer)
	GetTea(client io.Writer)
}
