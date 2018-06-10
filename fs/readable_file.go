package fs

import (
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
)

// An existing file that can be read or patched to an updated version
type ReadableFile struct {
	Filename string
}

func (readableFile *ReadableFile) Name() string {
	return "Existing file"
}

func (readableFile *ReadableFile) Get(client io.Writer, message http.RequestMessage) {
	readableFile.Head(client, message)
	slice := readableFile.makeSliceOfTargetFile(message)
	slice.WriteBody(client)
}

func (readableFile *ReadableFile) Head(client io.Writer, message http.RequestMessage) {
	slice := readableFile.makeSliceOfTargetFile(message)
	slice.WriteStatus(client)
	slice.WriteContentHeaders(client)
	msg.WriteEndOfMessageHeader(client)
}

func (readableFile *ReadableFile) makeSliceOfTargetFile(message http.RequestMessage) FileSlice {
	contentType := contentTypeFromFileExtension(readableFile.Filename)
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) != 1 {
		return &WholeFile{
			ContentType: contentType,
			Path:        readableFile.Filename,
		}
	}

	return ParseByteRange(rangeHeaders[0], readableFile.Filename, contentType)
}
