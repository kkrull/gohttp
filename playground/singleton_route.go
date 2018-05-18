package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewSingletonRoute(path string) *SingletonRoute {
	return &SingletonRoute{
		Singleton: &SingletonResource{Path: path},
	}
}

type SingletonRoute struct {
	Singleton *SingletonResource
}

func (route *SingletonRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.Singleton.Path {
		return nil
	}

	return requested.MakeResourceRequest(route.Singleton)
}

type SingletonResource struct {
	Path string
}

func (singleton *SingletonResource) Name() string {
	return "Singleton"
}

func (singleton *SingletonResource) Post(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.CreatedStatus)

	url := strings.Join([]string{singleton.Path, "data"}, "/")
	msg.WriteHeader(client, "Location", url)
	
	msg.WriteEndOfMessageHeader(client)
}
