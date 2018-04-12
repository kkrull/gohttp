package fs

import (
	"io"
	"mime"
	"os"
	"path"
	"strconv"

	"github.com/kkrull/gohttp/response"
)

type FileContents struct {
	Filename string
}

func (contents FileContents) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 200, "OK")
	contents.writeHeadersDescribingFile(client)
	response.WriteEndOfMessageHeader(client)

	file, _ := os.Open(contents.Filename)
	response.CopyToBody(client, file)
	return nil
}

func (contents FileContents) writeHeadersDescribingFile(client io.Writer) {
	response.WriteHeader(client, "Content-Type", contentTypeFromFileExtension(contents.Filename))
	info, _ := os.Stat(contents.Filename)
	response.WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}
