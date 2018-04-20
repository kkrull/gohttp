package fs

import "github.com/kkrull/gohttp/http"

func NewRoute(contentRootPath string) http.Route {
	return &fsRoute{
		ContentRootPath: contentRootPath,
		Controller:      &Controller{BaseDirectory: contentRootPath},
	}
}

type fsRoute struct {
	ContentRootPath string
	Controller      *Controller
}

func (route fsRoute) Route(requested *http.RequestLine) http.Request {
	switch requested.Method {
	case "GET":
		return &GetRequest{
			Controller: route.Controller,
			Target:     requested.Target,
		}
	case "HEAD":
		return &HeadRequest{
			Controller: route.Controller,
			Target:     requested.Target,
		}
	default:
		return nil
	}
}
