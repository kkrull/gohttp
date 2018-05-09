package http

import (
	"fmt"
	"io"
	"time"
)

// A null object for RequestLogger that does nothing
type noLogger struct{}

func (noLogger) Parsed(message RequestMessage) {}

type TextLogger struct {
	Writer io.Writer
}

func (logger TextLogger) Parsed(message RequestMessage) {
	fmt.Fprintf(logger.Writer, "%s : %s %s\n",
		time.Now().Format("2006-01-02 03:04:05 Z07:00"),
		message.Method(),
		message.Target())
	for _, header := range message.HeaderLines() {
		fmt.Fprintln(logger.Writer, header)
	}
}
