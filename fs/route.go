package fs

import "github.com/kkrull/gohttp/http"

type Route struct {
	ContentRootPath string
}

func (route Route) Route(method string, target string) http.Request {
	switch method {
	case "GET":
		return &GetRequest{
			BaseDirectory: route.ContentRootPath,
			Target:        target,
		}
	default:
		return nil
	}
}
