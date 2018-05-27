package playground

import (
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/success"
)

func NewCookieRoute(setTypePath, readTypePath string) *CookieRoute {
	return &CookieRoute{
		SetTypePath:  setTypePath,
		ReadTypePath: readTypePath,
	}
}

// Routes to register some content to return from the server in the form of a cookie
// and to ensure that the cookie sent back from the browser contains that content
type CookieRoute struct {
	SetTypePath string
	Registrar   *CookieRegistrar

	ReadTypePath string
}

func (route *CookieRoute) Route(requested http.RequestMessage) http.Request {
	switch requested.Path() {
	case route.SetTypePath:
		return requested.MakeResourceRequest(&CookieRegistrar{})
	default:
		return nil
	}
}

// Registers cookies
type CookieRegistrar struct{}

func (registrar *CookieRegistrar) Name() string {
	return "Cookie registrar"
}

func (registrar *CookieRegistrar) Get(client io.Writer, message http.RequestMessage) {
	cookieType := findQueryParameter(message, "type")
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteHeader(client, "Set-Cookie", cookieType)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("Eat a %s.", cookieType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.WriteBody(client, body)
}

func findQueryParameter(message http.RequestMessage, name string) string {
	for _, parameter := range message.QueryParameters() {
		if parameter.Name == name {
			return parameter.Value
		}
	}

	return ""
}
