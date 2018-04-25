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

func (route *ParameterRoute) Route(requested *http.RequestLine) http.Request {
	if requested.Target != "/parameters" {
		return nil
	}

	return http.MakeResourceRequest(requested, route.Decoder)
}

type ParameterDecoder interface {
	Name() string
	Get(client io.Writer, target string)
}

type TheDecoder struct{}

func (decoder *TheDecoder) Get(client io.Writer, target string) {
	panic("implement me")
}

func (decoder *TheDecoder) Name() string {
	panic("implement me")
}
