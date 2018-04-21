package teapot

import "io"

type GetRequest struct {
	Controller Controller
	Target string
}

func (request *GetRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

// Responds as a teapot that is aware of its own identity
type IdentityController struct {

}

func (controller *IdentityController) Get(client io.Writer, target string) {
}

type Controller interface {
	Get(client io.Writer, target string)
}
