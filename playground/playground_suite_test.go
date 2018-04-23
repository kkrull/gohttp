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
	getTarget     string
	headTarget    string
	optionsTarget string
	postTarget    string
	putTarget     string
}

func (mock *ControllerMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *ControllerMock) GetShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *ControllerMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *ControllerMock) HeadShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

func (mock *ControllerMock) Options(client io.Writer, target string) {
	mock.optionsTarget = target
}

func (mock *ControllerMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsTarget).To(Equal(target))
}

func (mock *ControllerMock) Post(client io.Writer, target string) {
	mock.postTarget = target
}

func (mock *ControllerMock) PostShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.postTarget).To(Equal(target))
}

func (mock *ControllerMock) Put(client io.Writer, target string) {
	mock.putTarget = target
}

func (mock *ControllerMock) PutShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.putTarget).To(Equal(target))
}
