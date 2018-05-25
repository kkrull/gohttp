package log_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "log")
}

/* WriterMock */

type WriterMock struct{}

func (mock *WriterMock) Length() int {
	return 0
}

func (mock *WriterMock) WriteLoggedRequests(client io.Writer) {}
