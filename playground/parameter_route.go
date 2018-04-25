package playground

import "github.com/kkrull/gohttp/http"

func NewParameterRoute() *ParameterRoute {
	return &ParameterRoute{}
}

type ParameterRoute struct{}

func (*ParameterRoute) Route(requested *http.RequestLine) http.Request {
	return nil
}
