package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

type AllowedMethodsController struct{}

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
