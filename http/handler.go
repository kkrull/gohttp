package http

import (
	"bufio"
	"io"

	"github.com/kkrull/gohttp/msg/servererror"
)

func NewConnectionHandler(router Router) ConnectionHandler {
	return &blockingConnectionHandler{Router: router}
}

// A ConnectionHandler that uses blocking I/O
type blockingConnectionHandler struct {
	Router Router
}

func (handler *blockingConnectionHandler) Handle(requestReader *bufio.Reader, responseWriter io.Writer) {
	request, routeErrorResponse := handler.Router.RouteRequest(requestReader)
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

func (handler *blockingConnectionHandler) Routes() []Route {
	return handler.Router.Routes()
}

type Router interface {
	RouteRequest(reader *bufio.Reader) (ok Request, err Response)
	Routes() []Route
}

type Request interface {
	Handle(client io.Writer) error
}

type Response interface {
	WriteTo(client io.Writer) error
	WriteHeader(client io.Writer) error
}
