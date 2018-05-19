package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/success"
)

func NewSingletonRoute(path string) *SingletonRoute {
	return &SingletonRoute{
		Singleton: &SingletonResource{CollectionPath: path},
	}
}

type SingletonRoute struct {
	Singleton *SingletonResource
}

func (route *SingletonRoute) Route(requested http.RequestMessage) http.Request {
	if !strings.HasPrefix(requested.Path(), route.Singleton.CollectionPath) {
		return nil
	}

	return requested.MakeResourceRequest(route.Singleton)
}

type SingletonResource struct {
	CollectionPath string
	data           []byte
}

func (singleton *SingletonResource) Name() string {
	return "Singleton"
}

func (singleton *SingletonResource) Get(client io.Writer, message http.RequestMessage) {
	if singleton.data == nil {
		clienterror.RespondNotFound(client, message.Path())
	} else if message.Path() != singleton.dataPath() {
		clienterror.RespondNotFound(client, message.Path())
	} else {
		msg.RespondInPlainText(client, success.OKStatus, singleton.data)
	}
}

func (singleton *SingletonResource) Post(client io.Writer, message http.RequestMessage) {
	singleton.data = message.Body()

	msg.WriteStatus(client, success.CreatedStatus)
	msg.WriteHeader(client, "Location", singleton.dataPath())
	msg.WriteEndOfMessageHeader(client)
}

func (singleton *SingletonResource) dataPath() string {
	return strings.Join([]string{singleton.CollectionPath, "data"}, "/")
}
