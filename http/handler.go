package http

import (
	"bufio"
	"io"

	"github.com/kkrull/gohttp/msg/servererror"
)

type ConnectionHandler struct {
	Router Router
}

func (handler *ConnectionHandler) Handle(requestReader *bufio.Reader, responseWriter io.Writer) {
	request, routeErrorResponse := handler.Router.ParseRequest(requestReader)
	if routeErrorResponse != nil {
		routeErrorResponse.WriteTo(responseWriter)
		return
	}

	requestError := request.Handle(responseWriter)
	if requestError != nil {
		response := servererror.InternalServerError{}
		response.WriteTo(responseWriter)
	}
}

type Router interface {
	AddRoute(route Route)
	ParseRequest(reader *bufio.Reader) (ok Request, routeError Response)
}

type Request interface {
	Handle(client io.Writer) error
}

type Response interface {
	WriteTo(client io.Writer) error
}
