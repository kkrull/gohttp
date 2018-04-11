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

func (response DirectoryListing) WriteTo(client io.Writer) {
	writeStatusLine(client, 200, "OK")
	writeHeader(client, "Content-Type", "text/plain")

	message := response.messageListingFiles()
	writeHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	writeEndOfMessageHeader(client)

	writeBody(client, message.String())
}

func (response DirectoryListing) messageListingFiles() *bytes.Buffer {
	message := &bytes.Buffer{}
	for _, file := range response.Files {
		fmt.Fprintf(message, "%s\n", file.Name())
	}

	return message
}

type FileContents struct {
	Filename string
}

func (response FileContents) WriteTo(client io.Writer) {
	writeStatusLine(client, 200, "OK")
	response.writeHeadersDescribingFile(client)
	writeEndOfMessageHeader(client)

	file, _ := os.Open(response.Filename)
	copyToBody(client, file)
}

func (response FileContents) writeHeadersDescribingFile(client io.Writer) {
	writeHeader(client, "Content-Type", contentTypeFromFileExtension(response.Filename))
	info, _ := os.Stat(response.Filename)
	writeHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}
