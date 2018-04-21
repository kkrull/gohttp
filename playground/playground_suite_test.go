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

type OptionControllerMock struct {
	optionsReceivedTarget string
}

func (mock *OptionControllerMock) Options(client io.Writer, target string) {
	mock.optionsReceivedTarget = target
}

func (mock *OptionControllerMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsReceivedTarget).To(Equal(target))
}
