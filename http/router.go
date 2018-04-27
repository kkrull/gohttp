package http

import (
	"bufio"
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

func (router RequestLineRouter) routeRequest(requested *requestMessage) (ok Request, notImplemented Response) {
	for _, route := range router.routes {
		request := route.Route(requested)
		if request != nil {
			return request, nil
		}
	}

	return nil, requested.NotImplemented()
}

type RequestParser interface {
	Parse(reader *bufio.Reader) (ok *requestMessage, err Response)
}

type RequestMessage interface {
	MakeResourceRequest(resource Resource) Request
	Method() string
	Target() string
	QueryParameters() []QueryParameter
}

type Route interface {
	Route(requested RequestMessage) Request
}

type QueryParameter struct {
	Name, Value string
}