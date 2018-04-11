package clientError

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kkrull/gohttp/response"
)

type BadRequest struct {
	DisplayText string
}

func (badRequest BadRequest) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 400, "Bad Request")
	return nil
}

type NotFound struct {
	Target string
}

func (notFound NotFound) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 404, "Not Found")
	response.WriteHeader(client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", notFound.Target)
	response.WriteHeader(client, "Content-Length", strconv.Itoa(len(message)))
	response.WriteEndOfMessageHeader(client)

	response.WriteBody(client, message)
	return nil
}
