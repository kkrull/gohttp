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

	notFoundReceived http.RequestMessage
	notFoundReturns  http.Resource
}

func (mock *ResourceFactoryMock) DirectoryListingResourceReturns(listing http.Resource) {
	mock.listingReturns = listing
}

func (mock *ResourceFactoryMock) DirectoryListingResource(message http.RequestMessage, files []string) http.Resource {
	mock.listingReceivedMessage = message
	mock.listingReceivedFiles = files
	return mock.listingReturns
}

func (mock *ResourceFactoryMock) DirectoryListingShouldHaveReceived(urlPath string, files []string) {
	ExpectWithOffset(1, mock.listingReceivedMessage).NotTo(BeNil())
	ExpectWithOffset(1, mock.listingReceivedMessage.Path()).To(Equal(urlPath))
	ExpectWithOffset(1, mock.listingReceivedFiles).To(Equal(files))
}

func (mock *ResourceFactoryMock) ExistingFileResourceReturns(file http.Resource) {
	mock.existingFileReturns = file
}

func (mock *ResourceFactoryMock) ExistingFileResource(message http.RequestMessage, path string) http.Resource {
	mock.existingFileReceivedMessage = message
	mock.existingFileReceivedPath = path
	return mock.existingFileReturns
}

func (mock *ResourceFactoryMock) ExistingFileShouldHaveReceived(urlPath string, fsPath string) {
	ExpectWithOffset(1, mock.existingFileReceivedMessage).NotTo(BeNil())
	ExpectWithOffset(1, mock.existingFileReceivedMessage.Path()).To(Equal(urlPath))
	ExpectWithOffset(1, mock.existingFileReceivedPath).To(Equal(fsPath))
}

func (mock *ResourceFactoryMock) NotFoundResource(message http.RequestMessage) http.Resource {
	mock.notFoundReceived = message
	return mock.notFoundReturns
}

func (mock *ResourceFactoryMock) NotFoundResourceReturns(notFound http.Resource) {
	mock.notFoundReturns = notFound
}

func (mock *ResourceFactoryMock) NotFoundShouldHaveReceived(path string) {
	ExpectWithOffset(1, mock.notFoundReceived).NotTo(BeNil())
	ExpectWithOffset(1, mock.notFoundReceived.Path()).To(Equal(path))
}
