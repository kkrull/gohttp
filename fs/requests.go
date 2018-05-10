package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

type ReadOnlyFileSystem struct {
	BaseDirectory string
}

func (controller *ReadOnlyFileSystem) Name() string {
	return "Readonly file system"
}

func (controller *ReadOnlyFileSystem) Get(client io.Writer, message http.RequestMessage) {
	response := controller.determineResponse(message.Path())
	response.WriteTo(client)
}

func (controller *ReadOnlyFileSystem) Head(client io.Writer, message http.RequestMessage) {
	response := controller.determineResponse(message.Path())
	response.WriteHeader(client)
}

func (controller *ReadOnlyFileSystem) determineResponse(requestedPath string) http.Response {
	resolvedPath := path.Join(controller.BaseDirectory, requestedPath)
	info, err := os.Stat(resolvedPath)
	if err != nil {
		return &clienterror.NotFound{Path: requestedPath}
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedPath)
		return &DirectoryListing{
			Files:      readFileNames(files),
			HrefPrefix: requestedPath}
	} else {
		return &ExistingFile{Filename: resolvedPath}
	}
}

func readFileNames(files []os.FileInfo) []string {
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	return fileNames
}
