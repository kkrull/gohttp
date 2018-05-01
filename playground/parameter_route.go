package playground

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
)

func NewParameterRoute() *ParameterRoute {
	return &ParameterRoute{Reporter: &AssignmentReporter{}}
}

type ParameterRoute struct {
	Reporter ParameterReporter
}

func (route *ParameterRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/parameters" {
		return nil
	}

	return requested.MakeResourceRequest(route.Reporter)
}

type ParameterReporter interface {
	Name() string
	Get(client io.Writer, req http.RequestMessage)
}

// Lists query parameters as simple assignment statements
type AssignmentReporter struct{}

func (reporter *AssignmentReporter) Name() string {
	return "Parameter Report"
}

func (reporter *AssignmentReporter) Get(client io.Writer, req http.RequestMessage) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentTypeHeader(client, "text/plain")

	body := reporter.makeBody(req)
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
