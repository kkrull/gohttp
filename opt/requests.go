package opt

import "io"

type OptionsRequest struct {
	Controller Controller
}

func (request *OptionsRequest) Handle(client io.Writer) error {
	request.Controller.Options(client)
	return nil
}

type Controller interface {
	Options(writer io.Writer)
}
