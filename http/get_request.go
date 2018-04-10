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
	Version       string
}

func (request *GetRequest) Handle(client *bufio.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		writeStatusLine(client, 404, "Not Found")
		writeHeader(client, "Content-Type", "text/plain")

		message := fmt.Sprintf("Not found: %s", request.Target)
		writeHeader(client, "Content-Length", strconv.Itoa(len(message)))
		writeEndOfHeader(client)

		writeBody(client, message)
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
		writeStatusLine(client, 200, "OK")
		writeHeader(client, "Content-Type", "text/plain")
		writeHeader(client, "Content-Length", strconv.FormatInt(info.Size(), 10))
		writeEndOfHeader(client)

		file, _ := os.Open(resolvedTarget)
		copyToBody(client, file)
	}

	client.Flush()
	return nil
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
