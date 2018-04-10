package http

import (
	"bufio"
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
	Version       string //TODO KDK: Not used
}

func (request *GetRequest) Handle(client *bufio.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		response := &NotFoundResponse{client: client}
		response.Issue(request.Target)
	} else if info.IsDir() {
		writeStatusLine(client, 200, "OK")
		writeHeader(client, "Content-Type", "text/plain")

		message := &bytes.Buffer{}
		messageWriter := bufio.NewWriter(message)
		files, _ := ioutil.ReadDir(resolvedTarget)
		for _, file := range files {
			messageWriter.WriteString(fmt.Sprintf("%s\n", file.Name()))
		}
		messageWriter.Flush()

		writeHeader(client, "Content-Length", strconv.Itoa(message.Len()))
		writeEndOfHeader(client)

		writeBody(client, message.String())
	} else {
		response := &FileContentsResponse{client: client}
		response.Issue(resolvedTarget)
	}

	client.Flush()
	return nil
}


type FileContentsResponse struct {
	client *bufio.Writer
}

func (response FileContentsResponse) Issue(filename string) {
	writeStatusLine(response.client, 200, "OK")
	writeHeader(response.client, "Content-Type", "text/plain")

	info, _ := os.Stat(filename)
	writeHeader(response.client, "Content-Length", strconv.FormatInt(info.Size(), 10))
	writeEndOfHeader(response.client)

	file, _ := os.Open(filename)
	copyToBody(response.client, file)
}


type NotFoundResponse struct {
	client *bufio.Writer
}

func (response NotFoundResponse) Issue(requestTarget string) {
	writeStatusLine(response.client, 404, "Not Found")
	writeHeader(response.client, "Content-Type", "text/plain")

	message := fmt.Sprintf("Not found: %s", requestTarget)
	writeHeader(response.client, "Content-Length", strconv.Itoa(len(message)))
	writeEndOfHeader(response.client)

	writeBody(response.client, message)
}


func writeStatusLine(client *bufio.Writer, status int, reason string) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", status, reason)
}

func writeHeader(client *bufio.Writer, name string, value string) {
	fmt.Fprintf(client, "%s: %s\r\n", name, value)
}

func writeEndOfHeader(client *bufio.Writer) {
	fmt.Fprint(client, "\r\n")
}

func copyToBody(client *bufio.Writer, bodyReader io.Reader) {
	io.Copy(client, bodyReader)
}

func writeBody(client *bufio.Writer, body string) {
	fmt.Fprint(client, body)
}
