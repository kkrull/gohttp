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
	getCoffeeCalled bool
	getTeaCalled    bool
}

func (mock *ControllerMock) GetCoffee(client io.Writer) {
	mock.getCoffeeCalled = true
}

func (mock *ControllerMock) GetCoffeeShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCoffeeCalled).To(BeTrue())
}

func (mock *ControllerMock) GetTea(client io.Writer) {
	mock.getTeaCalled = true
}

func (mock *ControllerMock) GetTeaShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getTeaCalled).To(BeTrue())
}
