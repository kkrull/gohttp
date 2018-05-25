package http

import (
	"bufio"
)

func NewRouter() *RequestLineRouter {
	return &RequestLineRouter{
		Parser: &LineRequestParser{},
		logger: noLogger{},
	}
}

// Routes requests based solely upon the first line in the request
type RequestLineRouter struct {
	Parser RequestParser
	logger RequestLogger
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

func (router *RequestLineRouter) LogRequests(logger RequestLogger) {
	router.logger = logger
}

func (router RequestLineRouter) RouteRequest(reader *bufio.Reader) (ok Request, err Response) {
	requested, err := router.Parser.Parse(reader)
	if err != nil {
		return nil, err
	}

	router.logger.Parsed(requested)
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

type RequestLogger interface {
	Parsed(message RequestMessage)
}

type RequestMessage interface {
	Method() string
	Target() string
	Version() string
	Path() string
	QueryParameters() []QueryParameter

	HeaderLines() []string
	HeaderValues(field string) (values []string)
	Body() []byte

	MakeResourceRequest(resource Resource) Request
}

type Route interface {
	Route(requested RequestMessage) Request
}

type QueryParameter struct {
	Name, Value string
}

// A null object for RequestLogger that does nothing
type noLogger struct{}

func (noLogger) Parsed(message RequestMessage) {}
