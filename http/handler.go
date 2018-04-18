package http

import (
	"bufio"
	"io"

	"github.com/kkrull/gohttp/msg/servererror"
)

// A ConnectionHandler that uses blocking I/O
type BlockingConnectionHandler struct {
	Router Router
}

func (handler *BlockingConnectionHandler) Handle(requestReader *bufio.Reader, responseWriter io.Writer) {
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
	ParseRequest(reader *bufio.Reader) (ok Request, routeError Response)
}

type Request interface {
	Handle(client io.Writer) error
}

type Response interface {
	WriteTo(client io.Writer) error
	WriteHeader(client io.Writer) error
}
