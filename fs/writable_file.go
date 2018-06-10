package fs

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
	"github.com/kkrull/gohttp/msg/success"
)

func NewWritableFileRoute(urlPath, contentRootPath string) *WritableFileRoute {
	resource := &WritableFile{Filename: path.Join(contentRootPath, urlPath)}
	return &WritableFileRoute{
		UrlPath:      urlPath,
		FileResource: resource,
	}
}

type WritableFileRoute struct {
	UrlPath      string
	FileResource *WritableFile
}

func (route *WritableFileRoute) Route(requested http.RequestMessage) http.Request {
	if requested.Path() != route.UrlPath {
		return nil
	}

	return requested.MakeResourceRequest(route.FileResource)
}

// An existing file with read/write access
type WritableFile struct {
	Filename string
}

func (writableFile *WritableFile) Name() string {
	return "Writable file"
}

func (writableFile *WritableFile) Get(client io.Writer, message http.RequestMessage) {
	writableFile.Head(client, message)
	slice := writableFile.makeSliceOfTargetFile(message)
	slice.WriteBody(client)
}

func (writableFile *WritableFile) Head(client io.Writer, message http.RequestMessage) {
	slice := writableFile.makeSliceOfTargetFile(message)
	slice.WriteStatus(client)
	slice.WriteContentHeaders(client)
	msg.WriteEndOfMessageHeader(client)
}

func (writableFile *WritableFile) makeSliceOfTargetFile(message http.RequestMessage) FileSlice {
	contentType := contentTypeFromFileExtension(writableFile.Filename)
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) != 1 {
		return &WholeFile{
			ContentType: contentType,
			Path:        writableFile.Filename,
		}
	}

	return ParseByteRange(rangeHeaders[0], writableFile.Filename, contentType)
}

//func contentTypeFromFileExtension(filename string) string {
//	extension := path.Ext(filename)
//	if extension == "" {
//		return "text/plain"
//	}
//
//	return mime.TypeByExtension(extension)
//}

func (writableFile *WritableFile) Patch(client io.Writer, message http.RequestMessage) {
	conditionalHeader, err := onlyConditionalHeader(message)
	if err != nil {
		msg.WriteStatus(client, clienterror.ConflictStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	if !writableFile.preconditionMatches(conditionalHeader) {
		msg.WriteStatus(client, clienterror.PreconditionFailedStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	if err := writableFile.overwriteFile(message.Body()); err != nil {
		msg.WriteStatus(client, servererror.InternalServerErrorStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	writableFile.successfulPatch(client, message.Path())
}

func (writableFile *WritableFile) Put(client io.Writer, message http.RequestMessage) {
	if err := writableFile.overwriteFile(message.Body()); err != nil {
		msg.WriteStatus(client, servererror.InternalServerErrorStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	writableFile.successfulPut(client, message.Path())
}

//func onlyConditionalHeader(message http.RequestMessage) (string, error) {
//	conditionalHeaders := message.HeaderValues("If-Match")
//	switch len(conditionalHeaders) {
//	case 0:
//		return "", &noConditionalHeadersError{}
//	case 1:
//		return conditionalHeaders[0], nil
//	default:
//		return "", &ambiguousConditionalHeadersError{}
//	}
//}

func (writableFile *WritableFile) validatorTag() string {
	return "\"" + writableFile.fileContentsHash() + "\""
}

func (writableFile *WritableFile) fileContentsHash() string {
	h := sha1.New()
	file, _ := os.Open(writableFile.Filename)
	defer file.Close()
	io.Copy(h, file)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//type noConditionalHeadersError struct{}
//
//func (noConditionalHeadersError) Error() string {
//	return "No If-Match header found"
//}
//
//type ambiguousConditionalHeadersError struct{}
//
//func (ambiguousConditionalHeadersError) Error() string {
//	return "Too many If-Match headers found"
//}

func (writableFile *WritableFile) preconditionMatches(preconditionHeader string) bool {
	currentETag := writableFile.fileContentsHash()
	return preconditionHeader == currentETag
}

func (writableFile *WritableFile) overwriteFile(body []byte) error {
	return ioutil.WriteFile(writableFile.Filename, body, os.ModePerm)
}

func (writableFile *WritableFile) successfulPatch(client io.Writer, path string) {
	msg.WriteStatus(client, success.NoContentStatus)
	msg.WriteHeader(client, "Content-Location", path)
	msg.WriteHeader(client, "ETag", writableFile.validatorTag())
	msg.WriteEndOfMessageHeader(client)
}

func (writableFile *WritableFile) successfulPut(client io.Writer, path string) {
	msg.WriteStatus(client, success.OKStatus)
	msg.WriteEndOfMessageHeader(client)
}

//// A view of all/part of a file
//type FileSlice interface {
//	WriteStatus(writer io.Writer)
//	WriteContentHeaders(writer io.Writer)
//	WriteBody(writer io.Writer)
//}
