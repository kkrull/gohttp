package teapot_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTeapot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "teapot")
}

type ControllerMock struct {
	getTarget string
}

func (mock *ControllerMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *ControllerMock) GetShouldHaveReceivedTarget(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}
