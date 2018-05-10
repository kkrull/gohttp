package fs_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "fs")
}

/* FileSystemResourceMock */

type FileSystemResourceMock struct {
	getPath    string
	headTarget string
}

func (mock *FileSystemResourceMock) Name() string {
	return "File system mock"
}

func (mock *FileSystemResourceMock) Get(client io.Writer, message http.RequestMessage) {
	mock.getPath = message.Path()
}

func (mock *FileSystemResourceMock) GetShouldHaveReceived(path string) {
	ExpectWithOffset(1, mock.getPath).To(Equal(path))
}

func (mock *FileSystemResourceMock) Head(client io.Writer, message http.RequestMessage) {
	mock.headTarget = message.Target()
}

func (mock *FileSystemResourceMock) HeadShouldHaveReceived(target string) {
	ExpectWithOffset(1, mock.headTarget).To(Equal(target))
}

/* ResourceFactoryMock */

type ResourceFactoryMock struct {
	existingFileReturns         http.Resource
	existingFileReceivedMessage http.RequestMessage
	existingFileReceivedPath    string

	listingReceivedFiles   []string
	listingReceivedMessage http.RequestMessage
	listingReturns         http.Resource

	nonExistingReceived http.RequestMessage
	nonExistingReturns  http.Resource
}

func (mock *ResourceFactoryMock) DirectoryListingResource(message http.RequestMessage, files []string) http.Resource {
	mock.listingReceivedMessage = message
	mock.listingReceivedFiles = files
	return mock.listingReturns
}

func (mock *ResourceFactoryMock) DirectoryListingResourceReturns(listing http.Resource) {
	mock.listingReturns = listing
}

func (mock *ResourceFactoryMock) DirectoryListingShouldHaveReceived(urlPath string, files []string) {
	ExpectWithOffset(1, mock.listingReceivedMessage).NotTo(BeNil())
	ExpectWithOffset(1, mock.listingReceivedMessage.Path()).To(Equal(urlPath))
	ExpectWithOffset(1, mock.listingReceivedFiles).To(Equal(files))
}

func (mock *ResourceFactoryMock) ExistingFileResource(message http.RequestMessage, path string) http.Resource {
	mock.existingFileReceivedMessage = message
	mock.existingFileReceivedPath = path
	return mock.existingFileReturns
}

func (mock *ResourceFactoryMock) ExistingFileResourceReturns(file http.Resource) {
	mock.existingFileReturns = file
}

func (mock *ResourceFactoryMock) ExistingFileShouldHaveReceived(urlPath string, fsPath string) {
	ExpectWithOffset(1, mock.existingFileReceivedMessage).NotTo(BeNil())
	ExpectWithOffset(1, mock.existingFileReceivedMessage.Path()).To(Equal(urlPath))
	ExpectWithOffset(1, mock.existingFileReceivedPath).To(Equal(fsPath))
}

func (mock *ResourceFactoryMock) NonExistingResource(message http.RequestMessage) http.Resource {
	mock.nonExistingReceived = message
	return mock.nonExistingReturns
}

func (mock *ResourceFactoryMock) NonExistingResourceReturns(notFound http.Resource) {
	mock.nonExistingReturns = notFound
}

func (mock *ResourceFactoryMock) NonExistingResourceShouldHaveReceived(path string) {
	ExpectWithOffset(1, mock.nonExistingReceived).NotTo(BeNil())
	ExpectWithOffset(1, mock.nonExistingReceived.Path()).To(Equal(path))
}

/* Helpers */

func makeEmptyTestDirectory(testName string, fileMode os.FileMode) string {
	testPath := path.Join(".test", testName)
	Expect(os.RemoveAll(testPath)).To(Succeed())
	Expect(os.MkdirAll(testPath, fileMode)).To(Succeed())
	return testPath
}

func createTextFile(filename string, contents string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	byteContents := bytes.NewBufferString(contents).Bytes()
	bytesWritten, err := file.Write(byteContents)
	if err != nil {
		return err
	} else if bytesWritten != len(byteContents) {
		return fmt.Errorf("expected to write %d bytes to %s, but only wrote %d", len(byteContents), filename, bytesWritten)
	}

	return nil
}
