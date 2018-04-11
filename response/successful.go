package response

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"strconv"
)

type DirectoryListing struct {
	Files []os.FileInfo
}

func (listing DirectoryListing) WriteTo(client io.Writer) error {
	WriteStatusLine(client, 200, "OK")
	WriteHeader(client, "Content-Type", "text/plain")

	message := listing.messageListingFiles()
	WriteHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	WriteEndOfMessageHeader(client)

	WriteBody(client, message.String())
	return nil
}

func (listing DirectoryListing) messageListingFiles() *bytes.Buffer {
	message := &bytes.Buffer{}
	for _, file := range listing.Files {
		fmt.Fprintf(message, "%s\n", file.Name())
	}

	return message
}

type FileContents struct {
	Filename string
}

func (contents FileContents) WriteTo(client io.Writer) error {
	WriteStatusLine(client, 200, "OK")
	contents.writeHeadersDescribingFile(client)
	WriteEndOfMessageHeader(client)

	file, _ := os.Open(contents.Filename)
	CopyToBody(client, file)
	return nil
}

func (contents FileContents) writeHeadersDescribingFile(client io.Writer) {
	WriteHeader(client, "Content-Type", contentTypeFromFileExtension(contents.Filename))
	info, _ := os.Stat(contents.Filename)
	WriteHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}
