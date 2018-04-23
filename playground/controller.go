package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

type AllowedMethodsController struct{}

func (controller *AllowedMethodsController) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *AllowedMethodsController) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *AllowedMethodsController) Post(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *AllowedMethodsController) Put(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}

func (controller *AllowedMethodsController) Options(client io.Writer, target string) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)

	allowedMethods := controller.methodsAllowedFor(target)
	msg.WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
}

func (controller *AllowedMethodsController) methodsAllowedFor(target string) []string {
	switch target {
	case "/method_options":
		return []string{"GET", "HEAD", "POST", "OPTIONS", "PUT"}
	case "/method_options2":
		return []string{"GET", "HEAD", "OPTIONS"}
	default:
		return nil
	}
}
