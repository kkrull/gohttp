package msg

import (
	"fmt"
	"io"
)

func WriteStatusLine(client io.Writer, status int, reason string) {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", status, reason)
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
