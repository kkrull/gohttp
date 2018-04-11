package response

import (
	"fmt"
	"io"
)

type BadRequest struct {
	DisplayText string
}

func (response BadRequest) WriteTo(client io.Writer) error {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 400, "Bad Request")
	return nil
}

type InternalServerError struct{}

func (response InternalServerError) WriteTo(client io.Writer) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 500, "Internal Server Error")
}

type NotImplemented struct {
	Method string
}

func (response NotImplemented) WriteTo(client io.Writer) error {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 501, "Not Implemented")
	return nil
}
