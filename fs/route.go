package fs

import "github.com/kkrull/gohttp/http"

func NewRoute(contentRootPath string) http.Route {
	return &route{
		ContentRootPath: contentRootPath,
		Controller:      &Controller{BaseDirectory: contentRootPath},
	}
}

type route struct {
	ContentRootPath string
	Controller      *Controller
}

func (route route) Route(requested *http.RequestLine) http.Request {
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
