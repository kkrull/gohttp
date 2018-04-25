package http

import (
	"bufio"

	"github.com/kkrull/gohttp/msg/servererror"
)

func NewRouter() *RequestLineRouter {
	return &RequestLineRouter{
		Parser: &LineRequestParser{},
	}
}

// Routes requests based solely upon the first line in the request
type RequestLineRouter struct {
	Parser RequestParser
	routes []Route
}

func (router *RequestLineRouter) AddRoute(route Route) {
	router.routes = append(router.routes, route)
}

func (router *RequestLineRouter) Routes() []Route {
	routes := make([]Route, len(router.routes))
	copy(routes, router.routes)
	return routes
}

func (router RequestLineRouter) RouteRequest(reader *bufio.Reader) (ok Request, err Response) {
	requested, err := router.Parser.Parse(reader)
	if err != nil {
		return nil, err
	}

	return router.routeRequest(requested)
}

func (router RequestLineRouter) routeRequest(requested *RequestLine) (ok Request, notImplemented Response) {
	for _, route := range router.routes {
		request := route.Route(requested)
		if request != nil {
			return request, nil
		}
	}

	return nil, requested.NotImplemented()
}

type RequestParser interface {
	Parse(reader *bufio.Reader) (ok *RequestLine, err Response)
}

type Route interface {
	Route(requested *RequestLine) Request
}

type RequestLine struct {
	Method          string
	Target          string
	QueryParameters map[string]string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}
