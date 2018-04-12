package fs

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"strconv"

	"github.com/kkrull/gohttp/response"
)

type DirectoryListing struct {
	Files []string
}

func (listing DirectoryListing) WriteTo(client io.Writer) error {
	response.WriteStatusLine(client, 200, "OK")
	response.WriteHeader(client, "Content-Type", "text/html")

	message := listing.messageListingFiles()
	response.WriteHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	response.WriteEndOfMessageHeader(client)

	response.WriteBody(client, message.String())
	return nil
}

func (listing DirectoryListing) messageListingFiles() *bytes.Buffer {
	message := &bytes.Buffer{}
	message.WriteString("<html>")
	message.WriteString("<head><title>gohttp</title></head>")

	message.WriteString("<body>")
	for _, file := range listing.Files {
		fmt.Fprintf(message, "<a href=\"/%s\">%s</a>", file, file)
	}

	message.WriteString("</body>")
	message.WriteString("</html>")
	return message
}

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
