package playground_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestPlayground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "playground")
}

/* Matchers */

func ShouldHaveNoBody(response *bytes.Buffer, status int, reason string) func() {
	return func() {
		responseMessage := httptest.ParseResponse(response)
		responseMessage.ShouldBeWellFormed()
		responseMessage.StatusShouldBe(status, reason)
		responseMessage.HeaderShould("Content-Length", Equal("0"))
		responseMessage.BodyShould(BeEmpty())
	}
}

func ContainSubstrings(values []string) types.GomegaMatcher {
	valueMatchers := make([]types.GomegaMatcher, len(values))
	for i, value := range values {
		valueMatchers[i] = ContainSubstring(value)
	}

	return SatisfyAll(valueMatchers...)
}

/* ReadOnlyResourceMock */

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

/* ReadWriteResourceMock */

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
