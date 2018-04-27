package fs_test

import (
	"io"
	"testing"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fs")
}

type FileSystemResourceMock struct {
	getPath    string
	headTarget string
}

func (mock *FileSystemResourceMock) Name() string {
	return "File system mock"
}

func (mock *FileSystemResourceMock) Get(client io.Writer, req http.RequestMessage) {
	mock.getPath = req.Path()
}

func (mock *FileSystemResourceMock) GetShouldHaveReceived(path string) {
	ExpectWithOffset(1, mock.getPath).To(Equal(path))
}

func (mock *FileSystemResourceMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *FileSystemResourceMock) HeadShouldHaveReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}
