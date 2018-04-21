package teapot

import "io"

type GetCoffeeRequest struct {

}

func (request *GetCoffeeRequest) Handle(client io.Writer) error {
	panic("implement me")
}
