package http

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type RFC7230RequestParser struct {
	BaseDirectory string
}

func (parser RFC7230RequestParser) ParseRequest(reader *bufio.Reader) (ok Request, parseError Response) {
	request, err := parser.parseRequestLine(reader)
	if err != nil {
		return nil, err
	}

	headerError := parseHeaderLines(reader)
	if headerError != nil {
		return nil, headerError
	}

	return request, nil
}

func (parser RFC7230RequestParser) parseRequestLine(reader *bufio.Reader) (*GetRequest, Response) {
	requestLine, err := readCRLFLine(reader)
	if err != nil {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	fields := strings.Split(requestLine, " ")
	if len(fields) != 3 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	switch fields[0] {
	case "GET":
		return &GetRequest{
			BaseDirectory: parser.BaseDirectory,
			Target:        fields[1],
		}, nil
	default:
		return nil, &ParseError{StatusCode: 501, Reason: "Not Implemented"}
	}
}

func parseHeaderLines(reader *bufio.Reader) *ParseError {
	isBlankLineBetweenHeadersAndBody := func(line string) bool { return line == "" }

	for {
		line, err := readCRLFLine(reader)
		if err != nil {
			return &ParseError{StatusCode: 400, Reason: "Bad Request"}
		} else if isBlankLineBetweenHeadersAndBody(line) {
			return nil
		}
	}
}

func readCRLFLine(reader *bufio.Reader) (string, error) {
	maybeEndsInCR, _ := reader.ReadString('\r')
	if len(maybeEndsInCR) == 0 {
		return "", MissingEndOfHeaderCRLF{}
	} else if !strings.HasSuffix(maybeEndsInCR, "\r") {
		return "", MalformedHeaderLine{}
	}

	nextCharacter, _ := reader.ReadByte()
	if nextCharacter != '\n' {
		return "", MalformedHeaderLine{}
	}

	trimmed := strings.TrimSuffix(maybeEndsInCR, "\r")
	return trimmed, nil
}

type MissingEndOfHeaderCRLF struct{}

func (MissingEndOfHeaderCRLF) Error() string {
	return "end of input before terminating CRLF"
}

type MalformedHeaderLine struct{}

func (MalformedHeaderLine) Error() string {
	return "line in request header not ending in CRLF"
}

type ParseError struct {
	StatusCode int
	Reason     string
}

func (parseError ParseError) WriteTo(client io.Writer) error {
	fmt.Fprintf(client, "HTTP/1.1 %d %s\r\n", parseError.StatusCode, parseError.Reason)
	return nil
}
