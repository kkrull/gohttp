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

type ReadOnlyResourceMock struct {
	getTarget     string
	headTarget    string
	optionsTarget string
	postTarget    string
	putTarget     string
}

func (mock *ReadOnlyResourceMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *ReadOnlyResourceMock) GetShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *ReadOnlyResourceMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *ReadOnlyResourceMock) HeadShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

func (mock *ReadOnlyResourceMock) Options(client io.Writer, target string) {
	mock.optionsTarget = target
}

func (mock *ReadOnlyResourceMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsTarget).To(Equal(target))
}

type ReadWriteResourceMock struct {
	getTarget     string
	headTarget    string
	optionsTarget string
	postTarget    string
	putTarget     string
}

func (mock *ReadWriteResourceMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *ReadWriteResourceMock) GetShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *ReadWriteResourceMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *ReadWriteResourceMock) HeadShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

func (mock *ReadWriteResourceMock) Options(client io.Writer, target string) {
	mock.optionsTarget = target
}

func (mock *ReadWriteResourceMock) OptionsShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.optionsTarget).To(Equal(target))
}

func (mock *ReadWriteResourceMock) Post(client io.Writer, target string) {
	mock.postTarget = target
}

func (mock *ReadWriteResourceMock) PostShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.postTarget).To(Equal(target))
}

func (mock *ReadWriteResourceMock) Put(client io.Writer, target string) {
	mock.putTarget = target
}

func (mock *ReadWriteResourceMock) PutShouldHaveBeenReceived(target string) {
	ExpectWithOffset(1, mock.putTarget).To(Equal(target))
}
