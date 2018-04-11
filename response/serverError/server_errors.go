package serverError

import (
	"io"

	"github.com/kkrull/gohttp/response"
)

type InternalServerError struct{}

func (internalError InternalServerError) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 500, "Internal Server Error")
	return nil
}

type NotImplemented struct {
	Method string
}

func (notImplemented NotImplemented) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 501, "Not Implemented")
	return nil
}
