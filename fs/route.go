package fs

import "github.com/kkrull/gohttp/http"

func NewRoute(contentRootPath string) http.Route {
	return &route{ContentRootPath: contentRootPath}
}

type route struct {
	ContentRootPath string
}

func (route route) Route(requested *http.RequestLine) http.Request {
	switch requested.Method {
	case "GET":
		return &GetRequest{
			BaseDirectory: route.ContentRootPath,
			Target:        requested.Target,
		}
	default:
		return nil
	}
}
