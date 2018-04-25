package http

import (
	"bufio"
)

type RequestLineRouter struct {
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
	parser := &RequestMessageParser{reader: reader}
	requested, err := parser.Parse()
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

type Route interface {
	Route(requested *RequestLine) Request
}
