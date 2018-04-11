package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"strconv"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
}

func (request *GetRequest) Handle(client io.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		response := &NotFoundResponse{Target: request.Target}
		response.WriteTo(client)
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		response := &DirectoryListingResponse{Files: files}
		response.WriteTo(client)
	} else {
		response := &FileContentsResponse{Filename: resolvedTarget}
		response.WriteTo(client)
	}

	return nil
}

type DirectoryListingResponse struct {
	Files []os.FileInfo
}

func (response DirectoryListingResponse) WriteTo(client io.Writer) {
	writeStatusLine(client, 200, "OK")
	writeHeader(client, "Content-Type", "text/plain")

	message := response.messageListingFiles()
	writeHeader(client, "Content-Length", strconv.Itoa(message.Len()))
	writeEndOfMessageHeader(client)

	writeBody(client, message.String())
}

func (response DirectoryListingResponse) messageListingFiles() *bytes.Buffer {
	message := &bytes.Buffer{}
	for _, file := range response.Files {
		fmt.Fprintf(message, "%s\n", file.Name())
	}

	return message
}

type FileContentsResponse struct {
	Filename string
}

func (response FileContentsResponse) WriteTo(client io.Writer) {
	writeStatusLine(client, 200, "OK")
	response.writeHeadersDescribingFile(client)
	writeEndOfMessageHeader(client)

	file, _ := os.Open(response.Filename)
	copyToBody(client, file)
}

func (response FileContentsResponse) writeHeadersDescribingFile(client io.Writer) {
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

type NotFoundResponse struct {
	Target string
}

func (response NotFoundResponse) WriteTo(client io.Writer) {
	writeStatusLine(client, 404, "Not Found")
	writeHeader(client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", response.Target)
	writeHeader(client, "Content-Length", strconv.Itoa(len(message)))
	writeEndOfMessageHeader(client)

	writeBody(client, message)
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
