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
	contentRange := existingFile.makeSliceOfTargetFile(message)
	contentRange.WriteBody(client)
}

func (existingFile *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	contentRange := existingFile.makeSliceOfTargetFile(message)
	contentRange.WriteStatus(client)
	contentRange.WriteContentSizeHeaders(client)
	msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))
	msg.WriteEndOfMessageHeader(client)
}

func (existingFile *ExistingFile) makeSliceOfTargetFile(message http.RequestMessage) FileSlice {
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) != 1 {
		return &WholeFile{path: existingFile.Filename}
	}

	contentRanges := ParseByteRangeSlices(rangeHeaders[0], existingFile.Filename)
	if len(contentRanges) != 1 {
		return &WholeFile{path: existingFile.Filename}
	}

	return contentRanges[0]
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}
