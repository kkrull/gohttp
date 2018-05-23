package log

import (
	"github.com/kkrull/gohttp/http"
)

func NewLogRoute(path string) http.Route {
	return &Route{
		Path:   path,
		Viewer: &Viewer{},
	}
}

type Route struct {
	Path   string
	Viewer *Viewer
}

func (route *Route) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Viewer)
}

// Views server logs
type Viewer struct{}

func (viewer *Viewer) Name() string {
	return "Log viewer"
}
