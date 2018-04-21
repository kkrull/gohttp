package opt

import "io"

type StaticCapabilitiesController struct {
}

func (controller *StaticCapabilitiesController) Options(writer io.Writer) {
	panic("implement me")
}

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
