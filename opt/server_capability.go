package opt

import "io"

type StaticCapabilityController struct {
}

func (controller *StaticCapabilityController) Options(writer io.Writer) {
	panic("implement me")
}
