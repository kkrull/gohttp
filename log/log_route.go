package log

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
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

func (viewer *Viewer) Get(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, clienterror.UnauthorizedStatus)
	msg.WriteHeader(client, "WWW-Authenticate", "Basic realm=\"logs\"")
	msg.WriteEndOfMessageHeader(client)
}
