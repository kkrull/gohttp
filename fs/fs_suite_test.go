package fs_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fs")
}

type FileSystemResourceMock struct {
	getTarget  string
	headTarget string
}

func (mock *FileSystemResourceMock) Get(client io.Writer, target string) {
	mock.getTarget = target
}

func (mock *FileSystemResourceMock) GetShouldHaveReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}

func (mock *FileSystemResourceMock) Head(client io.Writer, target string) {
	mock.headTarget = target
}

func (mock *FileSystemResourceMock) HeadShouldHaveReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}
