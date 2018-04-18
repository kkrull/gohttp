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
	Controller Controller
	Target     string
}

func (request *HeadRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}

type Controller struct {
	BaseDirectory string
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

func (controller *Controller) Head(client io.Writer, target string) {
	resolvedTarget := path.Join(controller.BaseDirectory, target)
	response := controller.determineResponse(target, resolvedTarget)
	response.WriteHeader(client)
}

func (controller *Controller) determineResponse(requested, resolved string) http.Response {
	info, err := os.Stat(resolved)
	if err != nil {
		return &clienterror.NotFound{Target: requested}
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolved)
		return &DirectoryListing{
			Files:      readFileNames(files),
			HrefPrefix: requested}
	} else {
		return &FileContents{Filename: resolved}
	}
}
