package playground

import (
	"io"

	"github.com/kkrull/gohttp/msg"
)

// Handles various read requests, but doesn't actually do anything
type ReadableNopResource struct{}

func (controller *ReadableNopResource) Name() string {
	return "Readonly NOP"
}

func (controller *ReadableNopResource) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *ReadableNopResource) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

// Handles various read/write requests, but doesn't actually do anything
type ReadWriteNopResource struct{}

func (controller *ReadWriteNopResource) Name() string {
	return "Read/Write NOP"
}

func (controller *ReadWriteNopResource) Get(client io.Writer, target string) {
	controller.Head(client, target)
}

func (controller *ReadWriteNopResource) Head(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Post(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func (controller *ReadWriteNopResource) Put(client io.Writer, target string) {
	writeOKWithNoBody(client)
}

func writeOKWithNoBody(client io.Writer) {
	msg.WriteStatusLine(client, 200, "OK")
	msg.WriteContentLengthHeader(client, 0)
	msg.WriteEndOfMessageHeader(client)
}
