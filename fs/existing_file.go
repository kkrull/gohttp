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
	currentETag := existingFile.fileContentsHash()
	conditionalHeaders := message.HeaderValues("If-Match")
	if len(conditionalHeaders) != 1 {
		msg.WriteStatus(client, clienterror.ConflictStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	conditionalHeader := conditionalHeaders[0]
	if conditionalHeader != currentETag {
		msg.WriteStatus(client, clienterror.PreconditionFailedStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	if err := ioutil.WriteFile(existingFile.Filename, message.Body(), os.ModePerm); err != nil {
		msg.WriteStatus(client, servererror.InternalServerErrorStatus)
		msg.WriteEndOfMessageHeader(client)
		return
	}

	msg.WriteStatus(client, success.NoContentStatus)
	msg.WriteHeader(client, "Content-Location", message.Path())
	msg.WriteHeader(client, "ETag", existingFile.validatorTag())
	msg.WriteEndOfMessageHeader(client)
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

// A view of all/part of a file
type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}
