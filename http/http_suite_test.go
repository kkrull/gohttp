package http_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "http")
}

/* Helpers */

func makeReader(template string, values ...interface{}) *bufio.Reader {
	text := fmt.Sprintf(template, values...)
	return bufio.NewReader(bytes.NewBufferString(text))
}
