package playground_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlayground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "playground")
}

type ControllerMock struct {
	optionsReceivedTarget string
}

func (mock *ControllerMock) GetShouldHaveBeenReceived(target string) {
}

func (mock *ControllerMock) Options(client io.Writer, target string) {
	mock.optionsReceivedTarget = target
}
func (mock *ControllerMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsReceivedTarget).To(Equal(target))
}
