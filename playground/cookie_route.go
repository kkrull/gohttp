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
		SetTypePath: setTypePath,
		Registrar:   &CookieRegistrar{},

		ReadTypePath: readTypePath,
		Monster:      &CookieMonster{},
	}
}

// Routes to register some content to return from the server in the form of a cookie
// and to ensure that the cookie sent back from the browser contains that content
type CookieRoute struct {
	SetTypePath string
	Registrar   *CookieRegistrar

	ReadTypePath string
	Monster      *CookieMonster
}

func (route *CookieRoute) Route(requested http.RequestMessage) http.Request {
	switch requested.Path() {
	case route.ReadTypePath:
		return requested.MakeResourceRequest(route.Monster)
	case route.SetTypePath:
		return requested.MakeResourceRequest(route.Registrar)
	default:
		return nil
	}
}

// "C" IS FOR COOKIE
type CookieMonster struct{}

func (monster *CookieMonster) Name() string {
	return "Cookie Monster"
}

func (monster *CookieMonster) Get(client io.Writer, message http.RequestMessage) {
	sessionCookie, err := singleHeader(message, "Cookie")
	if err != nil {
		monster.badCookieState(client)
	} else {
		cookieType := sessionCookie
		monster.preferredCookieState(client, cookieType)
	}
}

func (monster *CookieMonster) badCookieState(client io.Writer) {
	msg.WriteStatus(client, clienterror.BadRequestStatus)
	msg.WriteEndOfMessageHeader(client)
}

func (monster *CookieMonster) preferredCookieState(client io.Writer, cookieType string) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("mmmm %s", cookieType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.WriteBody(client, body)
}

// Registers cookies
type CookieRegistrar struct{}

func (registrar *CookieRegistrar) Name() string {
	return "Cookie registrar"
}

func (registrar *CookieRegistrar) Get(client io.Writer, message http.RequestMessage) {
	cookieType, err := singleQueryParameter(message, "type")
	if err != nil {
		registrar.invalidTypeState(client)
	} else {
		registrar.typeSetState(client, cookieType)
	}
}

func (registrar *CookieRegistrar) invalidTypeState(client io.Writer) {
	msg.WriteStatus(client, clienterror.BadRequestStatus)
	msg.WriteEndOfMessageHeader(client)
}

func (registrar *CookieRegistrar) typeSetState(client io.Writer, cookieType string) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteHeader(client, "Set-Cookie", cookieType)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("Eat a %s.", cookieType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.WriteBody(client, body)
}

func singleHeader(message http.RequestMessage, field string) (value string, err error) {
	values := message.HeaderValues(field)
	switch len(values) {
	case 0:
		return "", &missingValueError{what: field}
	case 1:
		return values[0], nil
	default:
		return "", &tooManyValuesError{what: field}
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
		return "", &missingValueError{what: "type parameter"}
	case 1:
		return values[0], nil
	default:
		return "", &tooManyValuesError{what: "type parameter"}
	}
}

// Resulting from a lack of any value for something that needed a value
type missingValueError struct {
	what string
}

func (err *missingValueError) Error() string {
	return fmt.Sprintf("no value provided for %s", err.what)
}

// Resulting from too many values for something that needed a single value
type tooManyValuesError struct {
	what string
}

func (err *tooManyValuesError) Error() string {
	return fmt.Sprintf("too many values for %s", err.what)
}
