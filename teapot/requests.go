package teapot

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type GetCoffeeRequest struct {
	Controller Controller
}

func (request *GetCoffeeRequest) Handle(client io.Writer) error {
	request.Controller.GetCoffee(client)
	return nil
}

type GetTeaRequest struct {
	Controller Controller
}

func (request *GetTeaRequest) Handle(client io.Writer) error {
	request.Controller.GetTea(client)
	return nil
}

// Responds as a teapot that is aware of its own identity
type IdentityController struct{}

func (controller *IdentityController) GetCoffee(client io.Writer) {
	body := "I'm a teapot"
	writeHeaders(client, body)
	msg.WriteBody(client, body)
}

func writeHeaders(client io.Writer, body string) {
	msg.WriteStatusLine(client, 418, "I'm a teapot")
	msg.WriteContentTypeHeader(client, "text/plain")
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
}

func (controller *IdentityController) GetTea(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
