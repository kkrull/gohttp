package response

import "io"

type InternalServerError struct{}

func (response InternalServerError) WriteTo(client io.Writer) {
	writeStatusLine(client, 500, "Internal Server Error")
}

type NotImplemented struct {
	Method string
}

func (response NotImplemented) WriteTo(client io.Writer) error {
	writeStatusLine(client, 501, "Not Implemented")
	return nil
}
