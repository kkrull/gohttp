package playground

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
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
	data []byte
}

func (singleton *SingletonResource) Name() string {
	return "Singleton"
}

func (singleton *SingletonResource) Get(client io.Writer, message http.RequestMessage) {
	if singleton.data == nil {
		msg.WriteStatus(client, clienterror.NotFoundStatus)
		msg.WriteContentTypeHeader(client, "text/plain")

		body := fmt.Sprintf("Not found: %s", singleton.dataUrl())
		msg.WriteContentLengthHeader(client, len(body))
		msg.WriteEndOfMessageHeader(client)

		msg.WriteBody(client, body)
	} else {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentTypeHeader(client, "text/plain")
		msg.WriteContentLengthHeader(client, len(singleton.data))
		msg.WriteEndOfMessageHeader(client)
		msg.CopyToBody(client, bytes.NewReader(singleton.data))
	}
}

func (singleton *SingletonResource) Post(client io.Writer, message http.RequestMessage) {
	singleton.data = message.Body()

	msg.WriteStatus(client, success.CreatedStatus)
	msg.WriteHeader(client, "Location", singleton.dataUrl())
	msg.WriteEndOfMessageHeader(client)
}

func (singleton *SingletonResource) dataUrl() string {
	return strings.Join([]string{singleton.Path, "data"}, "/")
}
