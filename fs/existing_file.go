package fs

import (
	"io"
	"mime"
	"os"
	"path"
	"strconv"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
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
	contentRange := existingFile.optionalByteRange(message)
	if contentRange == nil {
		existingFile.writeWholeFile(client)
	} else {
		contentRange.Copy(existingFile.Filename, client)
	}
}

func (existingFile *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	contentRange := existingFile.optionalByteRange(message)
	if contentRange == nil {
		msg.WriteStatus(client, success.OKStatus)
		info, _ := os.Stat(existingFile.Filename)
		msg.WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
	} else {
		msg.WriteStatus(client, success.PartialContentStatus)
		msg.WriteHeader(client, "Content-Length", strconv.Itoa(contentRange.Length()))
		msg.WriteHeader(client, "Content-Range", contentRange.ContentRange())
	}

	msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))
	msg.WriteEndOfMessageHeader(client)
}

func (existingFile *ExistingFile) optionalByteRange(message http.RequestMessage) *byteRange {
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) != 1 {
		return nil
	}

	info, _ := os.Stat(existingFile.Filename)
	contentRanges := ParseByteRanges(rangeHeaders[0], info.Size())
	if len(contentRanges) != 1 {
		return nil
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

func (existingFile *ExistingFile) writeWholeFile(client io.Writer) {
	file, _ := os.Open(existingFile.Filename)
	msg.CopyToBody(client, file)
}
