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

type BadRequest struct {
	DisplayText string
}

func (response BadRequest) WriteTo(client io.Writer) error {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 400, "Bad Request")
	return nil
}

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

type InternalServerError struct{}

func (response InternalServerError) WriteTo(client io.Writer) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 500, "Internal Server Error")
}

type NotFound struct {
	Target string
}

func (response NotFound) WriteTo(client io.Writer) {
	writeStatusLine(client, 404, "Not Found")
	writeHeader(client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", response.Target)
	writeHeader(client, "Content-Length", strconv.Itoa(len(message)))
	writeEndOfMessageHeader(client)

	writeBody(client, message)
}

type NotImplemented struct {
	Method string
}

func (response NotImplemented) WriteTo(client io.Writer) error {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", 501, "Not Implemented")
	return nil
}

func writeStatusLine(client io.Writer, status int, reason string) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", status, reason)
}

func writeHeader(client io.Writer, name string, value string) {
	fmt.Fprintf(client, "%s: %s\r\n", name, value)
}

func writeEndOfMessageHeader(client io.Writer) {
	fmt.Fprint(client, "\r\n")
}

func copyToBody(client io.Writer, bodyReader io.Reader) {
	io.Copy(client, bodyReader)
}

func writeBody(client io.Writer, body string) {
	fmt.Fprint(client, body)
}
