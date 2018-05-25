package log_test

import (
	"bytes"
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "log")
}

/* RequestBufferStub */

type RequestBufferStub struct {
	NumBytesReturns int
	WriteToWill     *bytes.Buffer
}

func (stub *RequestBufferStub) NumBytes() int {
	return stub.NumBytesReturns
}

func (stub *RequestBufferStub) WriteTo(client io.Writer) {
	stub.WriteToWill.WriteTo(client)
}
