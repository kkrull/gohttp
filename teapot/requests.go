package teapot

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

type GetCoffeeRequest struct {
	Controller Teapot
}

func (request *GetCoffeeRequest) Handle(client io.Writer) error {
	request.Controller.GetCoffee(client)
	return nil
}

type GetTeaRequest struct {
	Controller Teapot
}

func (request *GetTeaRequest) Handle(client io.Writer) error {
	request.Controller.GetTea(client)
	return nil
}

// Responds as a teapot that is aware of its own identity
type IdentityTeapot struct{}

func (controller *IdentityTeapot) Name() string {
	return "teapot"
}

func (controller *IdentityTeapot) Get(client io.Writer, target string) {
	switch target {
	case "/coffee":
		controller.GetCoffee(client)
	case "/tea":
		controller.GetTea(client)
	default:
		panic("unknown target")
	}
}

func (controller *IdentityTeapot) GetCoffee(client io.Writer) {
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

func (controller *IdentityTeapot) GetTea(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
