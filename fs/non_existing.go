package fs

import (
	"fmt"
	"io"

	"github.com/kkrull/gohttp/http"
	"github.com/kkrull/gohttp/msg"
	"github.com/kkrull/gohttp/msg/clienterror"
)

type NonExisting struct {
	Path string
	body string
}

func (nonExisting *NonExisting) Name() string {
	return "Non-existing file"
}

func (nonExisting *NonExisting) Get(client io.Writer, message http.RequestMessage) {
	nonExisting.Head(client, message)
	msg.WriteBody(client, nonExisting.body)
}

func (nonExisting *NonExisting) Head(client io.Writer, message http.RequestMessage) {
	msg.WriteStatus(client, clienterror.NotFoundStatus)
	msg.WriteContentTypeHeader(client, "text/plain")

	nonExisting.body = fmt.Sprintf("Not found: %s", nonExisting.Path)
	msg.WriteContentLengthHeader(client, len(nonExisting.body))
	msg.WriteEndOfMessageHeader(client)
}
