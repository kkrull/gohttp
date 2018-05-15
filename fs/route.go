package fs

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
)

func NewRoute(contentRootPath string) http.Route {
	return &FileSystemRoute{
		ContentRootPath: contentRootPath,
		Factory:         &LocalResources{},
	}
}

type FileSystemRoute struct {
	ContentRootPath string
	Factory         ResourceFactory
}

func (route FileSystemRoute) Route(requested http.RequestMessage) http.Request {
	resource := route.resolveResource(requested)
	return requested.MakeResourceRequest(resource)
}

func (route FileSystemRoute) resolveResource(requested http.RequestMessage) http.Resource {
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

func readFileNames(files []os.FileInfo) []string {
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	return fileNames
}

type ResourceFactory interface {
	DirectoryListingResource(message http.RequestMessage, files []string) http.Resource
	ExistingFileResource(message http.RequestMessage, path string) http.Resource
	NonExistingResource(message http.RequestMessage) http.Resource
}
