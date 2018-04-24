package fs

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg/clienterror"
)

func NewRoute(contentRootPath string) http.Route {
	return &FileSystemRoute{
		ContentRootPath: contentRootPath,
		Resource:        &ReadOnlyFileSystem{BaseDirectory: contentRootPath},
	}
}

type FileSystemRoute struct {
	ContentRootPath string
	Resource        FileSystemResource
}

func (route FileSystemRoute) Route(requested *http.RequestLine) http.Request {
	switch requested.Method {
	case "GET":
		return &GetRequest{
			Controller: route.Resource,
			Target:     requested.Target,
		}
	case "HEAD":
		return &HeadRequest{
			Controller: route.Resource,
			Target:     requested.Target,
		}
	default:
		return clienterror.MethodNotAllowed("GET", "HEAD")
	}
}

// Represents files and directories on the file system
type FileSystemResource interface {
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
}

type GetRequest struct {
	Controller FileSystemResource
	Target     string
}

func (request *GetRequest) Handle(client io.Writer) error {
	request.Controller.Get(client, request.Target)
	return nil
}

type HeadRequest struct {
	Controller FileSystemResource
	Target     string
}

func (request *HeadRequest) Handle(client io.Writer) error {
	request.Controller.Head(client, request.Target)
	return nil
}
