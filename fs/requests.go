package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

type GetRequest struct {
	Controller *Controller
	Target     string
}

func (request *GetRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

type HeadRequest struct {
	Controller *Controller
	Target     string
}

func (request *HeadRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}

type Controller struct {
	BaseDirectory string
}

func (controller *Controller) Get(client io.Writer, target string) {
	response := controller.determineResponse(target)
	response.WriteTo(client)
}

func (controller *Controller) Head(client io.Writer, target string) {
	response := controller.determineResponse(target)
	response.WriteHeader(client)
}

func (controller *Controller) determineResponse(requestedTarget string) http.Response {
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
