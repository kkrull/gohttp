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

type TeapotMock struct {
	getCoffeeCalled bool
	getTeaCalled    bool
}

func (mock *TeapotMock) Name() string {
	return "teapot mock"
}

func (mock *TeapotMock) Get(client io.Writer, target string) {
	panic("implement me")
}

func (mock *TeapotMock) GetCoffee(client io.Writer) {
	mock.getCoffeeCalled = true
}

func (mock *TeapotMock) GetCoffeeShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCoffeeCalled).To(BeTrue())
}

func (mock *TeapotMock) GetTea(client io.Writer) {
	mock.getTeaCalled = true
}

func (mock *TeapotMock) GetTeaShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getTeaCalled).To(BeTrue())
}
