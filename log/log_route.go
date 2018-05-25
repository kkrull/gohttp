package log

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/success"
)

func NewLogRoute(path string, requests RequestBuffer) http.Route {
	return &Route{
		Path:   path,
		Viewer: &Viewer{Requests: requests},
	}
}

// A route for viewing request logs
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

// Views logs of HTTP requests
type Viewer struct {
	Requests RequestBuffer
}

func (viewer *Viewer) Name() string {
	return "Log viewer"
}

func (viewer *Viewer) Get(client io.Writer, message http.RequestMessage) {
	machine := logWriterStateMachine{
		requests: viewer.Requests,
		client:   client,
	}
	machine.FindAuthorizationHeader(message)
}

// State machine to go through the workflow of parsing and validating authorization, before writing request logs
type logWriterStateMachine struct {
	requests RequestBuffer
	client   io.Writer
}

func (state *logWriterStateMachine) FindAuthorizationHeader(message http.RequestMessage) {
	authorizations := message.HeaderValues("Authorization")
	switch len(authorizations) {
	case 0:
		state.unauthorized()
	case 1:
		state.parseAuthorization(authorizations[0])
	default:
		state.ambiguouslyAuthorized()
	}
}

func (state *logWriterStateMachine) parseAuthorization(authorization string) {
	fields := strings.Split(authorization, " ")
	state.testAuthorization(fields[0], fields[1])
}

func (state *logWriterStateMachine) testAuthorization(method string, encodedCredentials string) {
	if method != "Basic" {
		state.forbidden()
	} else if encodedCredentials != "YWRtaW46aHVudGVyMg==" {
		state.forbidden()
	} else {
		state.authorized()
	}
}

func (state *logWriterStateMachine) ambiguouslyAuthorized() {
	msg.WriteStatus(state.client, clienterror.BadRequestStatus)
	msg.WriteEndOfMessageHeader(state.client)
}

func (state *logWriterStateMachine) authorized() {
	msg.WriteStatus(state.client, success.OKStatus)
	msg.WriteContentLengthHeader(state.client, state.requests.NumBytes())
	msg.WriteContentTypeHeader(state.client, "text/plain")
	msg.WriteEndOfMessageHeader(state.client)
	state.requests.WriteTo(state.client)
}

func (state *logWriterStateMachine) forbidden() {
	msg.WriteStatus(state.client, clienterror.ForbiddenStatus)
	msg.WriteEndOfMessageHeader(state.client)
}

func (state *logWriterStateMachine) unauthorized() {
	msg.WriteStatus(state.client, clienterror.UnauthorizedStatus)
	msg.WriteHeader(state.client, "WWW-Authenticate", "Basic realm=\"logs\"")
	msg.WriteEndOfMessageHeader(state.client)
}

// Writes HTTP requests
type RequestBuffer interface {
	NumBytes() int
	WriteTo(client io.Writer)
}
