package playground

import "github.com/kkrull/gohttp/http"

func NewSingletonRoute(path string) *SingletonRoute {
	return &SingletonRoute{Path: path}
}

type SingletonRoute struct {
	Path string
}

func (route *SingletonRoute) Route(requested http.RequestMessage) http.Request {
	return nil
}

