package playground

import "io"

type StatelessOptionController struct {

}

func (controller *StatelessOptionController) Options(client io.Writer, target string) {
	panic("implement me")
}

