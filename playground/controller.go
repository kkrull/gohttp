package playground

import (
	"io"
	"strings"

	"github.com/kkrull/gohttp/msg"
)

// Handles various read requests, but doesn't actually do anything
type ReadableNopResource struct{}

func (controller *ReadableNopResource) Get(client io.Writer) {
	controller.Head(client)
}

func (controller *ReadableNopResource) Head(client io.Writer) {
	writeOKWithNoBody(client)
}

func (controller *ReadableNopResource) Options(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)

	allowedMethods := []string{"GET", "HEAD", "OPTIONS"}
	msg.WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
}

// Handles various read/write requests, but doesn't actually do anything
type ReadWriteNopResource struct{}

func (controller *ReadWriteNopResource) Get(client io.Writer) {
	controller.Head(client)
}

func (controller *ReadWriteNopResource) Head(client io.Writer) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Options(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)

	allowedMethods := []string{"GET", "HEAD", "POST", "OPTIONS", "PUT"}
	msg.WriteHeader(client, "Allow", strings.Join(allowedMethods, ","))
	msg.WriteEndOfMessageHeader(client)
}

func (controller *ReadWriteNopResource) Post(client io.Writer) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Put(client io.Writer) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
