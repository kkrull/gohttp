package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
)

func NewRoute(contentRootPath string) http.Route {
	return &NewFileSystemRoute{
		ContentRootPath: contentRootPath,
		Factory:         &LocalResources{},
	}
}

type NewFileSystemRoute struct {
	ContentRootPath string
	Factory         ResourceFactory
}

func (route NewFileSystemRoute) Route(requested http.RequestMessage) http.Request {
	resource := route.resolveResource(requested)
	return requested.MakeResourceRequest(resource)
}

func (route NewFileSystemRoute) resolveResource(requested http.RequestMessage) http.Resource {
	resolvedPath := path.Join(route.ContentRootPath, requested.Path())
	info, err := os.Stat(resolvedPath)
	if err != nil {
		return route.Factory.NonExistingResource(requested)
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedPath)
		return route.Factory.DirectoryListingResource(requested, readFileNames(files))
	} else {
		return route.Factory.ExistingFileResource(requested, resolvedPath)
	}
}

type ResourceFactory interface {
	DirectoryListingResource(message http.RequestMessage, files []string) http.Resource
	ExistingFileResource(message http.RequestMessage, path string) http.Resource
	NonExistingResource(message http.RequestMessage) http.Resource
}

type FileSystemRoute struct {
	ContentRootPath string
	Resource        FileSystemResource
}

func (route FileSystemRoute) Route(requested http.RequestMessage) http.Request {
	return requested.MakeResourceRequest(route.Resource)
}

// Represents files and directories on the file system
type FileSystemResource interface {
	Name() string
	Get(client io.Writer, message http.RequestMessage)
	Head(client io.Writer, message http.RequestMessage)
}
