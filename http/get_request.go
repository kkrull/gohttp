package http

import (
	"bufio"
	"fmt"
)

type GetRequest struct {
	Method  string
	Target  string
	Version string
}

func (request *GetRequest) Handle(conn *bufio.Writer) error {
	switch request.Target {
	case "/":
		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
		fmt.Fprint(conn, "Content-Length: 5\r\n")
		fmt.Fprint(conn, "Content-Type: text/plain\r\n")
		fmt.Fprint(conn, "\r\n")
		fmt.Fprintf(conn, "hello")
		return nil
	default:
		fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n")
		return nil
	}
}
