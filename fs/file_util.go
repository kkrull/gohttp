package fs

import (
	"io"
	"mime"
	"path"

	"github.com/kkrull/gohttp/http"
)

func contentTypeFromFileExtension(filename string) string {
	extension := path.Ext(filename)
	if extension == "" {
		return "text/plain"
	}

	return mime.TypeByExtension(extension)
}

func onlyConditionalHeader(message http.RequestMessage) (string, error) {
	conditionalHeaders := message.HeaderValues("If-Match")
	switch len(conditionalHeaders) {
	case 0:
		return "", &noConditionalHeadersError{}
	case 1:
		return conditionalHeaders[0], nil
	default:
		return "", &ambiguousConditionalHeadersError{}
	}
}

type noConditionalHeadersError struct{}

func (noConditionalHeadersError) Error() string {
	return "No If-Match header found"
}

type ambiguousConditionalHeadersError struct{}

func (ambiguousConditionalHeadersError) Error() string {
	return "Too many If-Match headers found"
}

// A view of all/part of a file
type FileSlice interface {
	WriteStatus(writer io.Writer)
	WriteContentHeaders(writer io.Writer)
	WriteBody(writer io.Writer)
}
