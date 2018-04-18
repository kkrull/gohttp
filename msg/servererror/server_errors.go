// HTTP 5xx Server Error responses from RFC 7231, Section 6.6
package servererror

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type InternalServerError struct{}

func (internalError InternalServerError) WriteTo(client io.Writer) error {
	msg.WriteStatusLine(client, 500, "Internal Server Error")
	return nil
}

type NotImplemented struct {
	Method string
}

func (notImplemented NotImplemented) WriteTo(client io.Writer) error {
	msg.WriteStatusLine(client, 501, "Not Implemented")
	return nil
}
