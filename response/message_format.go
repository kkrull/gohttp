package response

import (
	"fmt"
	"io"
)

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
