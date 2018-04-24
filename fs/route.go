package fs

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
)

func NewRoute(contentRootPath string) http.Route {
	return &FileSystemRoute{
		ContentRootPath: contentRootPath,
		Resource:        &ReadOnlyFilesystem{BaseDirectory: contentRootPath},
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
		return &MethodNotAllowedRequest{SupportedMethods: []string{"GET", "HEAD"}}
	}
}

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

type MethodNotAllowedRequest struct {
	SupportedMethods []string
}

func (notAllowed *MethodNotAllowedRequest) Handle(client io.Writer) error {
	msg.WriteStatusLine(client, 405, "Method Not Allowed")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteHeader(client, "Allow", strings.Join(notAllowed.SupportedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
	return nil
}
