package fs

import (
	"io"
	"mime"
	"os"
	"path"
	"strconv"

	"github.com/kkrull/gohttp/msg"
)

type FileContents struct {
	Filename string
}

func (contents *FileContents) WriteTo(client io.Writer) error {
	contents.WriteHeader(client)
	contents.writeBody(client)
	return nil
}

func (contents *FileContents) WriteHeader(client io.Writer) error {
	msg.WriteStatusLine(client, 200, "OK")
	contents.writeHeadersDescribingFile(client)
	msg.WriteEndOfMessageHeader(client)
	return nil
}

func (contents FileContents) writeHeadersDescribingFile(client io.Writer) {
	msg.WriteHeader(client, "Content-Type", contentTypeFromFileExtension(contents.Filename))
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

func (contents *FileContents) writeBody(client io.Writer) {
	file, _ := os.Open(contents.Filename)
	msg.CopyToBody(client, file)
}
