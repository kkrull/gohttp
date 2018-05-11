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
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) == 1 {
		info, _ := os.Stat(existingFile.Filename)
		contentRanges := ParseByteRanges(rangeHeaders[0], info.Size())
		contentRanges[0].Copy(existingFile.Filename, client)
	} else {
		existingFile.writeWholeFile(client)
	}
}

func (existingFile *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	rangeHeaders := message.HeaderValues("Range")
	if len(rangeHeaders) == 1 {
		info, _ := os.Stat(existingFile.Filename)
		msg.WriteStatus(client, success.PartialContentStatus)
		msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))

		contentRanges := ParseByteRanges(rangeHeaders[0], info.Size())
		contentRange := contentRanges[0]
		msg.WriteHeader(client, "Content-Length", strconv.Itoa(contentRange.Length()))
		msg.WriteHeader(client, "Content-Range", contentRange.ContentRange())
		msg.WriteEndOfMessageHeader(client)
	} else {
		msg.WriteStatus(client, success.OKStatus)
		msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(existingFile.Filename))
		info, _ := os.Stat(existingFile.Filename)
		msg.WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
		msg.WriteEndOfMessageHeader(client)
	}
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
