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

func (contents *ExistingFile) Name() string {
	return "Existing file"
}

func (contents *ExistingFile) Get(client io.Writer, message http.RequestMessage) {
	contents.Head(client, message)
	contents.writeBody(client)
}

func (contents *ExistingFile) Head(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, success.OKStatus)
	contents.writeHeadersDescribingFile(client)
	msg.WriteEndOfMessageHeader(client)
}

func (contents ExistingFile) writeHeadersDescribingFile(client io.Writer) {
	msg.WriteContentTypeHeader(client, contentTypeFromFileExtension(contents.Filename))
	info, _ := os.Stat(contents.Filename)
	msg.WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}

func (contents *ExistingFile) writeBody(client io.Writer) {
	file, _ := os.Open(contents.Filename)
	msg.CopyToBody(client, file)
}
