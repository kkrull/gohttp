package mock

import (
	"bufio"
	"fmt"
	"io"

	. "github.com/onsi/gomega"
)

type Request struct {
	HandleReturns  string
	handleReceived io.Writer
}

func (mock *Request) Handle(writer io.Writer) error {
	mock.handleReceived = writer
	if mock.HandleReturns != "" {
		return fmt.Errorf(mock.HandleReturns)
	}

	return nil
}

func (mock *Request) VerifyHandle(writer *bufio.Writer) {
	ExpectWithOffset(1, mock.handleReceived).To(BeIdenticalTo(writer))
}
