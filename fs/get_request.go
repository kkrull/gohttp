package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/response/clienterror"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
}

func (request *GetRequest) Handle(client io.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		notFound := &clienterror.NotFound{Target: request.Target}
		notFound.WriteTo(client)
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		fileNames := make([]string, len(files))
		for i, file := range files {
			fileNames[i] = file.Name()
		}

		directoryListing := &DirectoryListing{Files: fileNames}
		directoryListing.WriteTo(client)
	} else {
		fileContents := &FileContents{Filename: resolvedTarget}
		fileContents.WriteTo(client)
	}

	return nil
}
