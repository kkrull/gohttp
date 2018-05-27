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
	ledger := &MemoryCookieLedger{}
	return &CookieRoute{
		SetTypePath: setTypePath,
		Registrar:   &CookieRegistrar{Ledger: ledger},

		ReadTypePath: readTypePath,
		Monster:      &CookieMonster{Ledger: ledger},
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
type CookieMonster struct {
	Ledger CookieLedger
}

func (monster *CookieMonster) Name() string {
	return "Cookie Monster"
}

func (monster *CookieMonster) Get(client io.Writer, message http.RequestMessage) {
	cookieType, _ := monster.Ledger.PreferredType()
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	body := fmt.Sprintf("mmmm %s", cookieType)
	msg.WriteContentLengthHeader(client, len(body))
	msg.WriteEndOfMessageHeader(client)
	msg.WriteBody(client, body)
}

// Registers cookies
type CookieRegistrar struct {
	Ledger CookieLedger
}

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

// Stores information about cookies in memory
type MemoryCookieLedger struct{}

func (ledger *MemoryCookieLedger) PreferredType() (string, error) {
	return "", fmt.Errorf("no preference has been defined")
}

// Keeps track of the cookie monster's preferred type of cookie
type CookieLedger interface {
	PreferredType() (string, error)
}
