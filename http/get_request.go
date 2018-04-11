package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
		response := &NotFoundResponse{client: client}
		response.IssueForTarget(request.Target)
	} else if info.IsDir() {
		files, _ := ioutil.ReadDir(resolvedTarget)
		response := &DirectoryListingResponse{client: client}
		response.IssueForFiles(files)
	} else {
		response := &TextFileContentsResponse{client: client}
		response.IssueForFile(resolvedTarget)
	}

	return nil
}

type DirectoryListingResponse struct {
	client io.Writer
}

func (response DirectoryListingResponse) IssueForFiles(files []os.FileInfo) {
	writeStatusLine(response.client, 200, "OK")
	writeHeader(response.client, "Content-Type", "text/plain")

	message := messageListingFiles(files)
	writeHeader(response.client, "Content-Length", strconv.Itoa(message.Len()))
	writeEndOfMessageHeader(response.client)

	writeBody(response.client, message.String())
}

func messageListingFiles(files []os.FileInfo) *bytes.Buffer {
	message := &bytes.Buffer{}
	for _, file := range files {
		fmt.Fprintf(message, "%s\n", file.Name())
	}

	return message
}

type TextFileContentsResponse struct {
	client io.Writer
}

func (response TextFileContentsResponse) IssueForFile(filename string) {
	writeStatusLine(response.client, 200, "OK")
	writeHeadersDescribingFile(response, filename)
	writeEndOfMessageHeader(response.client)

	file, _ := os.Open(filename)
	copyToBody(response.client, file)
}

func writeHeadersDescribingFile(response TextFileContentsResponse, filename string) {
	writeHeader(response.client, "Content-Type", "text/plain")
	info, _ := os.Stat(filename)
	writeHeader(response.client, "Content-Length", strconv.FormatInt(info.Size(), 10))
}

type NotFoundResponse struct {
	client io.Writer
}

func (response NotFoundResponse) IssueForTarget(requestTarget string) {
	writeStatusLine(response.client, 404, "Not Found")
	writeHeader(response.client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", requestTarget)
	writeHeader(response.client, "Content-Length", strconv.Itoa(len(message)))
	writeEndOfMessageHeader(response.client)

	writeBody(response.client, message)
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
