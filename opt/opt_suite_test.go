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

type ControllerMock struct {
	optionsCalled bool
}

func (mock *ControllerMock) Options(writer io.Writer) {
	mock.optionsCalled = true
}

func (mock *ControllerMock) OptionsShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.optionsCalled).To(BeTrue())
}
