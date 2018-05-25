package log

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/success"
)

func NewLogRoute(path string, requests Writer) http.Route {
	return &Route{
		Path: path,
		Viewer: &Viewer{
			RequestLog: requests,
		},
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
type Viewer struct {
	RequestLog Writer
}

func (viewer *Viewer) Name() string {
	return "Log viewer"
}

func (viewer *Viewer) Get(client io.Writer, message http.RequestMessage) {
	authorizations := message.HeaderValues("Authorization")
	if len(authorizations) == 0 {
		msg.WriteStatus(client, clienterror.UnauthorizedStatus)
		msg.WriteHeader(client, "WWW-Authenticate", "Basic realm=\"logs\"")
		msg.WriteEndOfMessageHeader(client)
		return
	}

	firstAuthorization := authorizations[0]
	fields := strings.Split(firstAuthorization, " ")

	method := fields[0]
	if method != "Basic" {
		msg.WriteStatus(client, clienterror.ForbiddenStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	encodedCredentials := fields[1]
	if encodedCredentials != "YWRtaW46aHVudGVyMg==" {
		msg.WriteStatus(client, clienterror.ForbiddenStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentLengthHeader(client, viewer.RequestLog.Length())
	msg.WriteContentTypeHeader(client, "text/plain")
	msg.WriteEndOfMessageHeader(client)
	viewer.RequestLog.WriteLoggedRequests(client)
}

type Writer interface {
	Length() int
	WriteLoggedRequests(client io.Writer)
}
