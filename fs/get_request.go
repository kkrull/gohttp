package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

type HeadRequest struct {
	BaseDirectory string
	Target        string
}

func (request *HeadRequest) Handle(client io.Writer) error {
	return nil
}

func (request *HeadRequest) Respond() http.Response {
	return nil
}

type GetRequest struct {
	BaseDirectory string
	Target        string
}

func (request *GetRequest) Handle(client io.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	response := request.determineResponse(resolvedTarget)
	response.WriteTo(client)
	return nil
}

func (request *GetRequest) Respond() http.Response {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	return request.determineResponse(resolvedTarget)
}

func (request *GetRequest) determineResponse(resolvedTarget string) http.Response {
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		return &clienterror.NotFound{Target: request.Target}
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		return &DirectoryListing{
			Files:      readFileNames(files),
			HrefPrefix: request.Target}
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
