package teapot

import (
	"io"
	"strconv"

	"github.com/kkrull/gohttp/msg"
)

type GetRequest struct {
	Controller Controller
	Target     string
}

func (request *GetRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

// Responds as a teapot that is aware of its own identity
type IdentityController struct {
	body string
}

func (controller *IdentityController) Get(client io.Writer, target string) {
	msg.WriteStatusLine(client, 418, "I'm a teapot")
	msg.WriteHeader(client, "Content-Type", "text/plain")
	body := "I'm a teapot"
	msg.WriteHeader(client, "Content-Length", strconv.Itoa(len(body)))
	msg.WriteEndOfMessageHeader(client)

	msg.WriteBody(client, body)
}

type Controller interface {
	Get(client io.Writer, target string)
}
