package playground

import "github.com/kkrull/gohttp/http"

func NewCookieRoute(setTypePath, readTypePath string) *CookieRoute {
	return &CookieRoute{
		SetTypePath:  setTypePath,
		ReadTypePath: readTypePath,
	}
}

type CookieRoute struct {
	SetTypePath  string
	ReadTypePath string
}

func (route *CookieRoute) Route(requested http.RequestMessage) http.Request {
	return nil
}
