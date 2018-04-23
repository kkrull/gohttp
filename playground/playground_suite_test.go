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

type ReadableControllerMock struct {
	getTarget     string
	headTarget    string
	optionsTarget string
	postTarget    string
	putTarget     string
}

func (mock *ReadableControllerMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *ReadableControllerMock) GetShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *ReadableControllerMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *ReadableControllerMock) HeadShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

func (mock *ReadableControllerMock) Options(client io.Writer, target string) {
	mock.optionsTarget = target
}

func (mock *ReadableControllerMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsTarget).To(Equal(target))
}

type WritableControllerMock struct {
	getTarget     string
	headTarget    string
	optionsTarget string
	postTarget    string
	putTarget     string
}

func (mock *WritableControllerMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *WritableControllerMock) GetShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *WritableControllerMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *WritableControllerMock) HeadShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

func (mock *WritableControllerMock) Options(client io.Writer, target string) {
	mock.optionsTarget = target
}

func (mock *WritableControllerMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsTarget).To(Equal(target))
}

func (mock *WritableControllerMock) Post(client io.Writer, target string) {
	mock.postTarget = target
}

func (mock *WritableControllerMock) PostShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.postTarget).To(Equal(target))
}

func (mock *WritableControllerMock) Put(client io.Writer, target string) {
	mock.putTarget = target
}

func (mock *WritableControllerMock) PutShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.putTarget).To(Equal(target))
}
