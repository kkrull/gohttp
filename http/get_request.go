package http

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/response"
	"github.com/kkrull/gohttp/response/clientError"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
}

func (request *GetRequest) Handle(client io.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		notFound := &clientError.NotFound{Target: request.Target}
		notFound.WriteTo(client)
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		directoryListing := &response.DirectoryListing{Files: files}
		directoryListing.WriteTo(client)
	} else {
		fileContents := &response.FileContents{Filename: resolvedTarget}
		fileContents.WriteTo(client)
	}

	return nil
}
