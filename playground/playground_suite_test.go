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
	getCalled     bool
	headCalled    bool
	optionsCalled bool
}

func (mock *ReadOnlyResourceMock) Get(client io.Writer) {
	mock.getCalled = true
}

func (mock *ReadOnlyResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadOnlyResourceMock) Head(client io.Writer) {
	mock.headCalled = true
}

func (mock *ReadOnlyResourceMock) HeadShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.headCalled).To(BeTrue())
}

func (mock *ReadOnlyResourceMock) Options(client io.Writer) {
	mock.optionsCalled = true
}

func (mock *ReadOnlyResourceMock) OptionsShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.optionsCalled).To(BeTrue())
}

type ReadWriteResourceMock struct {
	getCalled     bool
	headCalled    bool
	optionsCalled bool
	postCalled    bool
	putCalled     bool
}

func (mock *ReadWriteResourceMock) Get(client io.Writer) {
	mock.getCalled = true
}

func (mock *ReadWriteResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Head(client io.Writer) {
	mock.headCalled = true
}

func (mock *ReadWriteResourceMock) HeadShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.headCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Options(client io.Writer) {
	mock.optionsCalled = true
}

func (mock *ReadWriteResourceMock) OptionsShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.optionsCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Post(client io.Writer) {
	mock.postCalled = true
}

func (mock *ReadWriteResourceMock) PostShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.postCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Put(client io.Writer) {
	mock.putCalled = true
}

func (mock *ReadWriteResourceMock) PutShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.putCalled).To(BeTrue())
}
