package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

// Handles various read/write requests, but doesn't actually do anything
type WritableNopController struct{}

func (controller *WritableNopController) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *WritableNopController) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *WritableNopController) Post(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *WritableNopController) Put(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}

func (controller *WritableNopController) Options(client io.Writer, target string) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)

	allowedMethods := controller.methodsAllowedFor(target)
	msg.WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
}

func (controller *WritableNopController) methodsAllowedFor(target string) []string {
	switch target {
	case "/method_options":
		return []string{"GET", "HEAD", "POST", "OPTIONS", "PUT"}
	case "/method_options2":
		return []string{"GET", "HEAD", "OPTIONS"}
	default:
		return nil
	}
}
