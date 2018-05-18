package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewFormRoute() *FormRoute {
	return &FormRoute{
		Form: &SingletonForm{},
	}
}

type FormRoute struct {
	Form *SingletonForm
}

func (route *FormRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/form" {
		return nil
	}

	return requested.MakeResourceRequest(route.Form)
}

type SingletonForm struct{}

func (*SingletonForm) Name() string {
	return "Singleton form"
}

func (*SingletonForm) Post(client io.Writer, message http.RequestMessage) {
}
