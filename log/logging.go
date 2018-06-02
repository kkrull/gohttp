package log

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/kkrull/gohttp/http"
)

func NewBufferedRequestLogger() *TextLogger {
	return &TextLogger{buffer: bytes.NewBuffer(make([]byte, 256*5000))}
}

// Logs HTTP requests to a buffer in plain text
type TextLogger struct {
	buffer *bytes.Buffer
}

func (logger TextLogger) Parsed(message http.RequestMessage) {
	fmt.Fprintf(logger.buffer, "%s : %s %s %s\n",
		time.Now().Format("2006-01-02 03:04:05 Z07:00"),
		message.Method(),
		message.Target(),
		message.Version())
	//for _, header := range message.HeaderLines() {
	//	fmt.Fprintln(logger.buffer, header)
	//}
	//fmt.Fprintln(logger.buffer)
	//fmt.Fprintf(logger.buffer, "%s", message.Body())
}

func (logger TextLogger) NumBytes() int {
	return logger.buffer.Len()
}

func (logger TextLogger) WriteTo(client io.Writer) {
	logger.buffer.WriteTo(client)
}
