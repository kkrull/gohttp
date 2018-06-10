package playground

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewParameterRoute(path string) *ParameterRoute {
	return &ParameterRoute{
		Path:     path,
		Reporter: &AssignmentReporter{},
	}
}

type ParameterRoute struct {
	Path     string
	Reporter ParameterReporter
}

func (route *ParameterRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Reporter)
}

type ParameterReporter interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
}

// Lists query parameters as simple assignment statements
type AssignmentReporter struct{}

func (reporter *AssignmentReporter) Name() string {
	return "Parameter Report"
}

func (reporter *AssignmentReporter) Get(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := reporter.makeBody(message)
	msg.WriteContentLengthHeader(client, body.Len())
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, body.String())
}

func (reporter *AssignmentReporter) makeBody(requestMessage http.RequestMessage) *bytes.Buffer {
	body := &bytes.Buffer{}
	for _, parameter := range requestMessage.QueryParameters() {
		fmt.Fprintf(body, "%s = %s\n", parameter.Name, parameter.Value)
	}

	return body
}
