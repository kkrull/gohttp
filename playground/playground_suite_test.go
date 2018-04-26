package playground_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlayground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "playground")
}

/* ParameterDecoderMock */

type ParameterDecoderMock struct {
	getParameters map[string]string
}

func (mock *ParameterDecoderMock) Name() string {
	return "Parameter Decoder Mock"
}

func (mock *ParameterDecoderMock) Get(client io.Writer, target string) {
	//TODO KDK: Work here to get the parameters and expand the interface
	//mock.getParameters = map[string]string {
	//	"two": "2",
	//	"one": "1",
	//}
}

func (mock *ParameterDecoderMock) GetShouldHaveReceived(parameters map[string]string) {
	ExpectWithOffset(1, mock.getParameters).To(Equal(parameters))
}

/* ReadOnlyResourceMock */

type ReadOnlyResourceMock struct {
	getCalled  bool
	headCalled bool
}

func (mock *ReadOnlyResourceMock) Name() string {
	return "Readonly Mock"
}

func (mock *ReadOnlyResourceMock) Get(client io.Writer, target string) {
	mock.getCalled = true
}

func (mock *ReadOnlyResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadOnlyResourceMock) Head(client io.Writer, target string) {
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

func (mock *ReadWriteResourceMock) Get(client io.Writer, target string) {
	mock.getCalled = true
}

func (mock *ReadWriteResourceMock) GetShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.getCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Head(client io.Writer, target string) {
	mock.headCalled = true
}

func (mock *ReadWriteResourceMock) HeadShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.headCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Post(client io.Writer, target string) {
	mock.postCalled = true
}

func (mock *ReadWriteResourceMock) PostShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.postCalled).To(BeTrue())
}

func (mock *ReadWriteResourceMock) Put(client io.Writer, target string) {
	mock.putCalled = true
}

func (mock *ReadWriteResourceMock) PutShouldHaveBeenCalled() {
	ExpectWithOffset(1, mock.putCalled).To(BeTrue())
}

/* Helpers */

func handleRequest(router http.Route, method, target string) {
	requested := &http.RequestLine{TheMethod: method, TheTarget: target}
	routedRequest := router.Route(requested)
	ExpectWithOffset(1, routedRequest).NotTo(BeNil())

	routedRequest.Handle(&bytes.Buffer{})
}
