package servererror

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type InternalServerError struct{}

func (internalError *InternalServerError) WriteTo(client io.Writer) error {
	return internalError.WriteHeader(client)
}

func (internalError *InternalServerError) WriteHeader(client io.Writer) error {
	msg.WriteStatus(client, InternalServerErrorStatus)
	//msg.WriteEndOfMessageHeader(client)
	return nil
}

type NotImplemented struct {
	Method string
}

func (notImplemented *NotImplemented) WriteTo(client io.Writer) error {
	return notImplemented.WriteHeader(client)
}

func (notImplemented *NotImplemented) WriteHeader(client io.Writer) error {
	msg.WriteStatus(client, NotImplementedStatus)
	//msg.WriteEndOfMessageHeader(client)
	return nil
}
