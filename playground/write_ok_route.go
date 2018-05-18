package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewWriteOKRoute(path string) *WriteOKRoute {
	return &WriteOKRoute{
		Path: path,
		Form: &WriteOKResource{},
	}
}

type WriteOKRoute struct {
	Path string
	Form *WriteOKResource
}

func (route *WriteOKRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Form)
}

// A resource where it's OK to write
type WriteOKResource struct{}

func (*WriteOKResource) Name() string {
	return "Write OK"
}

func (resource *WriteOKResource) Post(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteEndOfMessageHeader(client)
}

func (resource *WriteOKResource) Put(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteEndOfMessageHeader(client)
}
