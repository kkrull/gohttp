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
	msg.RespondWithAllowHeader(client, success.OKStatus, singleton.allowedMethods(message.Path()))
}

func (singleton *SingletonResource) allowedMethods(path string) []string {
	switch path {
	case singleton.CollectionPath:
		return collectionMethods
	case singleton.dataPath():
		return dataMethods
	default:
		return []string{http.OPTIONS}
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
