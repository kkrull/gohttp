package fs

import (
	"io"

	"github.com/kkrull/gohttp/http"
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
	return http.MakeResourceRequest(requested, route.Resource)
}

// Represents files and directories on the file system
type FileSystemResource interface {
	Name() string
	Get(client io.Writer, target string)
	Head(client io.Writer, target string)
}
