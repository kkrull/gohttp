package teapot

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
)

// Responds as a teapot that is aware of its own identity
type IdentityTeapot struct{}

func (teapot *IdentityTeapot) Name() string {
	return "teapot"
}

func (teapot *IdentityTeapot) RespondsTo(path string) bool {
	switch path {
	case "/coffee", "/tea":
		return true
	default:
		return false
	}
}

func (teapot *IdentityTeapot) Get(client io.Writer, req http.RequestMessage) {
	var beverageRequestHandlers = map[string]func(writer io.Writer){
		"/coffee": teapot.getCoffee,
		"/tea":    teapot.getTea,
	}

	handler := beverageRequestHandlers[req.Path()]
	handler(client)
}

func (teapot *IdentityTeapot) getCoffee(client io.Writer) {
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

func (teapot *IdentityTeapot) getTea(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
