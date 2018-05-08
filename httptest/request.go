package httptest

import (
	"bufio"
	"fmt"
	"io"

	. "github.com/onsi/gomega"
)

type RequestMock struct {
	HandleReturns  string
	handleReceived io.Writer
}

func (mock *RequestMock) Handle(writer io.Writer) error {
	mock.handleReceived = writer
	if mock.HandleReturns != "" {
		return fmt.Errorf(mock.HandleReturns)
	}

	return nil
}

func (mock *RequestMock) VerifyHandle(writer *bufio.Writer) {
	ExpectWithOffset(1, mock.handleReceived).To(BeIdenticalTo(writer))
}
