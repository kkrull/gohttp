package http

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
	Version       string
}

func (request *GetRequest) Handle(client *bufio.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		writeStatusLine(client, 404, "Not Found")
	} else {
		writeStatusLine(client, 200, "OK")
		writeHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
		writeHeader(client, "Content-Type", "text/plain")
		writeEndOfHeader(client)

		file, _ := os.Open(resolvedTarget)
		writeBody(client, file)
	}

	client.Flush()
	return nil
}

func writeBody(client *bufio.Writer, body io.Reader) {
	io.Copy(client, body)
}

func writeStatusLine(client *bufio.Writer, status int, reason string) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", status, reason)
}

func writeHeader(client *bufio.Writer, name string, value string) {
	fmt.Fprintf(client, "%s: %s\r\n", name, value)
}

func writeEndOfHeader(client *bufio.Writer) {
	fmt.Fprintf(client, "\r\n")
}
