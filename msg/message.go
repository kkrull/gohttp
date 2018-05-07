package msg

import (
	"fmt"
	"io"
	"strconv"
)

func WriteStatus(client io.Writer, status Status) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", status.Code, status.Reason)
}

type Status struct {
	Code   uint
	Reason string
}

func WriteContentLengthHeader(client io.Writer, numBytes int) {
	WriteHeader(client, "Content-Length", strconv.Itoa(numBytes))
}

func WriteContentTypeHeader(client io.Writer, value string) {
	WriteHeader(client, "Content-Type", value)
}

func WriteHeader(client io.Writer, name string, value string) {
	fmt.Fprintf(client, "%s: %s\r\n", name, value)
}

func WriteEndOfMessageHeader(client io.Writer) {
	fmt.Fprint(client, "\r\n")
}

func CopyToBody(client io.Writer, bodyReader io.Reader) {
	io.Copy(client, bodyReader)
}

func WriteBody(client io.Writer, body string) {
	fmt.Fprint(client, body)
}
