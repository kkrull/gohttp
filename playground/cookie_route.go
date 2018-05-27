package playground

import (
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
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
	cookieType, err := singleQueryParameter(message, "type")
	if err != nil {
		invalidTypeState(client)
	} else {
		typeSetState(client, cookieType)
	}
}

func singleQueryParameter(message http.RequestMessage, name string) (value string, err error) {
	values := make([]string, 0)
	for _, parameter := range message.QueryParameters() {
		if parameter.Name == name {
			values = append(values, parameter.Value)
		}
	}

	switch len(values) {
	case 0:
		return "", fmt.Errorf("missing type parameter")
	case 1:
		return values[0], nil
	default:
		return "", fmt.Errorf("too many type parameters")
	}
}

func invalidTypeState(client io.Writer) {
	msg.WriteStatus(client, clienterror.BadRequestStatus)
	msg.WriteEndOfMessageHeader(client)
}

func typeSetState(client io.Writer, cookieType string) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteHeader(client, "Set-Cookie", cookieType)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("Eat a %s.", cookieType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.WriteBody(client, body)
}
