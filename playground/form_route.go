package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewFormRoute(path string) *FormRoute {
	return &FormRoute{
		Path: path,
		Form: &SingletonForm{},
	}
}

type FormRoute struct {
	Path string
	Form *SingletonForm
}

func (route *FormRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Form)
}

type SingletonForm struct{}

func (*SingletonForm) Name() string {
	return "Singleton form"
}

func (form *SingletonForm) Post(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteEndOfMessageHeader(client)
}

func (form *SingletonForm) Put(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteEndOfMessageHeader(client)
}
