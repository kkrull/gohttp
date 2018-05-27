package playground_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/httptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlayground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "playground")
}

/* ParameterReporterMock */

type ParameterReporterMock struct{}

func (mock *ParameterReporterMock) Name() string {
	return "Parameter Reporter Mock"
}

func (mock *ParameterReporterMock) Get(client io.Writer, message http.RequestMessage) {}

/* ReadOnlyResourceMock */

type ReadOnlyResourceMock struct {
	getCalled  bool
	headCalled bool
}

func (mock *ReadOnlyResourceMock) Name() string {
	return "Readonly Mock"
}

func (mock *ReadOnlyResourceMock) Get(client io.Writer, message http.RequestMessage) {
	mock.getCalled = true
}

func (mock *ReadOnlyResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadOnlyResourceMock) Head(client io.Writer, message http.RequestMessage) {
	mock.headCalled = true
}

func (mock *ReadOnlyResourceMock) HeadShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.headCalled).To(BeTrue())
}

/* ReadWriteResourceMock */

type ReadWriteResourceMock struct {
	getCalled  bool
	headCalled bool
	postCalled bool
	putCalled  bool
}

func (mock *ReadWriteResourceMock) Name() string {
	return "Read/Write Mock"
}

func (mock *ReadWriteResourceMock) Get(client io.Writer, message http.RequestMessage) {
	mock.getCalled = true
}

func (mock *ReadWriteResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Head(client io.Writer, message http.RequestMessage) {
	mock.headCalled = true
}

func (mock *ReadWriteResourceMock) HeadShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.headCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Post(client io.Writer, message http.RequestMessage) {
	mock.postCalled = true
}

func (mock *ReadWriteResourceMock) PostShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.postCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Put(client io.Writer, message http.RequestMessage) {
	mock.putCalled = true
}

func (mock *ReadWriteResourceMock) PutShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.putCalled).To(BeTrue())
}

/* Helpers */

func invokeResourceMethod(invokeMethod httpResourceMethod, request http.RequestMessage) *httptest.ResponseMessage {
	response := &bytes.Buffer{}
	invokeMethod(response, request)
	return httptest.ParseResponse(response)
}

type httpResourceMethod = func(io.Writer, http.RequestMessage)

func handleRequest(router http.Route, method, path string) {
	requested := http.NewRequestMessage(method, path)
	routedRequest := router.Route(requested)
	ExpectWithOffset(1, routedRequest).NotTo(BeNil())

	routedRequest.Handle(&bytes.Buffer{})
}
