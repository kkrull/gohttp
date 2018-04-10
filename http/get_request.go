package http

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
)

type GetRequest struct {
	BaseDirectory string
	Target        string
	Version       string
}

func (request *GetRequest) Handle(response *bufio.Writer) error {
	resolvedTarget := path.Join(request.BaseDirectory, request.Target)
	info, err := os.Stat(resolvedTarget)
	if err != nil {
		fmt.Fprint(response, "HTTP/1.1 404 Not Found\r\n")
	} else {
		fmt.Fprintf(response, "HTTP/1.1 200 OK\r\n")
		fmt.Fprintf(response, "Content-Length: %d\r\n", info.Size())
		fmt.Fprintf(response, "Content-Type: text/plain\r\n")
		fmt.Fprintf(response, "\r\n")

		file, _ := os.Open(resolvedTarget)
		io.Copy(response, file)
	}

	response.Flush()
	return nil
}
