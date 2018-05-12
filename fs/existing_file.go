package fs

import (
	"io"
	"mime"
	"path"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
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
			Path: existingFile.Filename,
		}
	}

	slice := ParseByteRangeSlice(rangeHeaders[0], existingFile.Filename, contentType)
	if slice == nil {
		return &WholeFile{
			ContentType: contentType,
			Path: existingFile.Filename,
		}
	}

	return slice
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}