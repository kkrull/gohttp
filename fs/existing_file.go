package fs

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
	"github.com/kkrull/gohttp/msg/success"
)

type ExistingFile struct {
	Filename string
}

func (existingFile *ExistingFile) Name() string {
	return "Existing file"
}

func (existingFile *ExistingFile) Get(client io.Writer, message http.RequestMessage) {
	existingFile.Head(client, message)
	slice := existingFile.makeSliceOfTargetFile(message)
	slice.WriteBody(client)
}

func (existingFile *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	slice := existingFile.makeSliceOfTargetFile(message)
	slice.WriteStatus(client)
	slice.WriteContentHeaders(client)
	msg.WriteEndOfMessageHeader(client)
}

func (existingFile *ExistingFile) makeSliceOfTargetFile(message http.RequestMessage) FileSlice {
	contentType := contentTypeFromFileExtension(existingFile.Filename)
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) != 1 {
		return &WholeFile{
			ContentType: contentType,
			Path:        existingFile.Filename,
		}
	}

	return ParseByteRange(rangeHeaders[0], existingFile.Filename, contentType)
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}

func (existingFile *ExistingFile) Patch(client io.Writer, message http.RequestMessage) {
	conditionalHeader, err := onlyConditionalHeader(message)
	if err != nil {
		msg.WriteStatus(client, clienterror.ConflictStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	if !existingFile.preconditionMatches(conditionalHeader) {
		msg.WriteStatus(client, clienterror.PreconditionFailedStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	if err := existingFile.overwriteFile(message.Body()); err != nil {
		msg.WriteStatus(client, servererror.InternalServerErrorStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	existingFile.successfulPatch(client, message.Path())
}

func onlyConditionalHeader(message http.RequestMessage) (string, error) {
	conditionalHeaders := message.HeaderValues("If-Match")
	switch len(conditionalHeaders) {
	case 0:
		return "", &noConditionalHeadersError{}
	case 1:
		return conditionalHeaders[0], nil
	default:
		return "", &ambiguousConditionalHeadersError{}
	}
}

func (existingFile *ExistingFile) validatorTag() string {
	return "\"" + existingFile.fileContentsHash() + "\""
}

func (existingFile *ExistingFile) fileContentsHash() string {
	h := sha1.New()
	file, _ := os.Open(existingFile.Filename)
	defer file.Close()
	io.Copy(h, file)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type noConditionalHeadersError struct{}

func (noConditionalHeadersError) Error() string {
	return "No If-Match header found"
}

type ambiguousConditionalHeadersError struct{}

func (ambiguousConditionalHeadersError) Error() string {
	return "Too many If-Match headers found"
}

func (existingFile *ExistingFile) preconditionMatches(preconditionHeader string) bool {
	currentETag := existingFile.fileContentsHash()
	return preconditionHeader == currentETag
}

func (existingFile *ExistingFile) overwriteFile(body []byte) error {
	return ioutil.WriteFile(existingFile.Filename, body, os.ModePerm)
}

func (existingFile *ExistingFile) successfulPatch(client io.Writer, path string) {
	msg.WriteStatus(client, success.NoContentStatus)
	msg.WriteHeader(client, "Content-Location", path)
	msg.WriteHeader(client, "ETag", existingFile.validatorTag())
	msg.WriteEndOfMessageHeader(client)
}

// A view of all/part of a file
type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}
