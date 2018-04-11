package response

import (
	"fmt"
	"io"
	"strconv"
)

type BadRequest struct {
	DisplayText string
}

func (response BadRequest) WriteTo(client io.Writer) error {
	writeStatusLine(client, 400, "Bad Request")
	return nil
}

type NotFound struct {
	Target string
}

func (response NotFound) WriteTo(client io.Writer) error {
	writeStatusLine(client, 404, "Not Found")
	writeHeader(client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", response.Target)
	writeHeader(client, "Content-Length", strconv.Itoa(len(message)))
	writeEndOfMessageHeader(client)

	writeBody(client, message)
	return nil
}
