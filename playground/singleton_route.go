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

func (singleton *SingletonResource) Delete(client io.Writer, message http.RequestMessage) {
	if !singleton.isRequestForData(message) {
		clienterror.RespondMethodNotAllowed(client, collectionMethods)
	} else if singleton.hasData() {
		singleton.deleteData()
		success.RespondOkWithoutBody(client)
	} else {
		clienterror.RespondNotFound(client, message.Path())
	}
}

func (singleton *SingletonResource) Get(client io.Writer, message http.RequestMessage) {
	if singleton.hasData() && singleton.isRequestForData(message) {
		success.RespondOKWithKnownBody(client, "text/plain", singleton.data)
	} else {
		clienterror.RespondNotFound(client, message.Path())
	}
}

func (singleton *SingletonResource) Options(client io.Writer, message http.RequestMessage) {
	if message.Path() == singleton.CollectionPath {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentLengthHeader(client, 0)
		msg.WriteHeader(client, "Allow", strings.Join(collectionMethods, ","))
		msg.WriteEndOfMessageHeader(client)
		return
	} else if singleton.isRequestForData(message) {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentLengthHeader(client, 0)
		msg.WriteHeader(client, "Allow", strings.Join(dataMethods, ","))
		msg.WriteEndOfMessageHeader(client)
		return
	} else {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentLengthHeader(client, 0)
		msg.WriteHeader(client, "Allow", strings.Join([]string{http.OPTIONS}, ","))
		msg.WriteEndOfMessageHeader(client)
		return
	}
}

func (singleton *SingletonResource) Post(client io.Writer, message http.RequestMessage) {
	singleton.setData(message.Body())
	msg.WriteStatus(client, success.CreatedStatus)
	msg.WriteHeader(client, "Location", singleton.dataPath())
	msg.WriteEndOfMessageHeader(client)
}

func (singleton *SingletonResource) Put(client io.Writer, message http.RequestMessage) {
	singleton.setData(message.Body())
	success.RespondOkWithoutBody(client)
}

func (singleton *SingletonResource) deleteData() {
	singleton.data = nil
}

func (singleton *SingletonResource) hasData() bool {
	return singleton.data != nil
}

func (singleton *SingletonResource) setData(body []byte) {
	singleton.data = body
}

func (singleton *SingletonResource) isRequestForData(message http.RequestMessage) bool {
	return message.Path() == singleton.dataPath()
}

func (singleton *SingletonResource) dataPath() string {
	return strings.Join([]string{singleton.CollectionPath, "data"}, "/")
}

var (
	collectionMethods = []string{http.OPTIONS, http.POST}
	dataMethods       = []string{http.OPTIONS, http.DELETE, http.GET, http.PUT}
)
