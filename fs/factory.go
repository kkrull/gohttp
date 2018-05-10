package fs

import "github.com/kkrull/gohttp/http"

// Factory for Resources in this package
type LocalResources struct{}

func (*LocalResources) DirectoryListingResource(message http.RequestMessage, files []string) http.Resource {
	return &DirectoryListing{
		Files:      files,
		HrefPrefix: message.Path(),
	}
}

func (*LocalResources) ExistingFileResource(message http.RequestMessage, path string) http.Resource {
	return &ExistingFile{Filename: path}
}

func (*LocalResources) NonExistingResource(message http.RequestMessage) http.Resource {
	return &NonExisting{Path: message.Path()}
}
