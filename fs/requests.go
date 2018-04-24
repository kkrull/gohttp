package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

type ReadOnlyFilesystem struct {
	BaseDirectory string
}

func (controller *ReadOnlyFilesystem) Get(client io.Writer, target string) {
	response := controller.determineResponse(target)
	response.WriteTo(client)
}

func (controller *ReadOnlyFilesystem) Head(client io.Writer, target string) {
	response := controller.determineResponse(target)
	response.WriteHeader(client)
}

func (controller *ReadOnlyFilesystem) determineResponse(requestedTarget string) http.Response {
	resolvedTarget := path.Join(controller.BaseDirectory, requestedTarget)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		return &clienterror.NotFound{Target: requestedTarget}
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		return &DirectoryListing{
			Files:      readFileNames(files),
			HrefPrefix: requestedTarget}
	} else {
		return &FileContents{Filename: resolvedTarget}
	}
}

func readFileNames(files []os.FileInfo) []string {
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	return fileNames
}
