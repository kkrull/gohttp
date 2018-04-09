package http

import (
	"bufio"
	"strings"
)

type RFC7230RequestParser struct{}

func (parser RFC7230RequestParser) ParseRequest(reader *bufio.Reader) (Request, *ParseError) {
	request, err := parseRequestLine(reader)
	if err != nil {
		return nil, err
	}

	headerError := parseHeaderLines(reader)
	if headerError != nil {
		return nil, headerError
	}

	return request, nil
}

func parseRequestLine(reader *bufio.Reader) (*GetRequest, *ParseError) {
	requestLine, err := readCRLFLine(reader)
	if err != nil {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	fields := strings.Split(requestLine, " ")
	if len(fields) != 3 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	return &GetRequest{
		Method:  fields[0],
		Target:  fields[1],
		Version: fields[2],
	}, nil
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
