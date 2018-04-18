package fs

import "github.com/kkrull/gohttp/http"

type Route struct {
	ContentRootPath string
}

func (route Route) Route(requested *http.RequestLine) http.Request {
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
