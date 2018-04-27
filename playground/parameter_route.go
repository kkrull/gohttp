package playground

import (
	"io"

	"github.com/kkrull/gohttp/http"
)

func NewParameterRoute() *ParameterRoute {
	return &ParameterRoute{Decoder: &TheDecoder{}}
}

type ParameterRoute struct {
	Decoder ParameterDecoder
}

func (route *ParameterRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != "/parameters" {
		return nil
	}

	return requested.MakeResourceRequest(route.Decoder)
}

type ParameterDecoder interface {
	Name() string
	Get(client io.Writer, req http.RequestMessage)
}

type TheDecoder struct{}

func (decoder *TheDecoder) Name() string {
	return "Parameters"
}

func (decoder *TheDecoder) Get(client io.Writer, req http.RequestMessage) {
	panic("implement me")
}
