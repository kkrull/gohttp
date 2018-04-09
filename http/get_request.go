package http

import (
	"bufio"
	"fmt"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
	Version       string
}

func (request *GetRequest) Handle(response *bufio.Writer) error {
	switch request.Target {
	case "/":
		fmt.Fprint(response, "HTTP/1.1 200 OK\r\n")
		fmt.Fprint(response, "Content-Length: 5\r\n")
		fmt.Fprint(response, "Content-Type: text/plain\r\n")
		fmt.Fprint(response, "\r\n")
		fmt.Fprintf(response, "hello")
	default:
		fmt.Fprint(response, "HTTP/1.1 404 Not Found\r\n")
	}

	response.Flush()
	return nil
}
