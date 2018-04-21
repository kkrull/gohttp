package opt_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "opt")
}

type ServerCapabilityControllerMock struct {
	optionsCalled bool
}

func (mock *ServerCapabilityControllerMock) Options(writer io.Writer) {
	mock.optionsCalled = true
}

func (mock *ServerCapabilityControllerMock) OptionsShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.optionsCalled).To(BeTrue())
}
