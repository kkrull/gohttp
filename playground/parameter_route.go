package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
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

// Lists parameters as simple assignment statements
type AssignmentReporter struct{}

func (decoder *AssignmentReporter) Name() string {
	return "Parameters"
}

func (decoder *AssignmentReporter) Get(client io.Writer, req http.RequestMessage) {
	panic("implement me")
}
