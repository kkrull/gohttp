package http

import (
	"bufio"
	"strings"
	"fmt"
)

const maxLengthOfFieldInRequestLine = 8000

type RequestParser interface {
	ParseRequest(reader *bufio.Reader) (*Request, *ParseError)
}

type RFC7230RequestParser struct{}

func (parser RFC7230RequestParser) ParseRequest(reader *bufio.Reader) (*Request, *ParseError) {
	requestLine, err := readCRLFLine(reader)
	if err != nil {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	fields := strings.Split(requestLine, " ")
	if len(fields) != 3 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	request := &Request{
		Method: fields[0],
		Target: fields[1],
		Version: fields[2],
	}

	for {
		line, err := readCRLFLine(reader)
		if err != nil {
			return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
		} else if line == "" {
			return request, nil
		}
	}
}

func readCRLFLine(reader *bufio.Reader) (string, error) {
	maybeEndsInCR, _ := reader.ReadString('\r')
	if len(maybeEndsInCR) == 0 {
		return "", fmt.Errorf("end of input")
	} else if !strings.HasSuffix(maybeEndsInCR, "\r") {
		return "", fmt.Errorf("request line not ending in CR")
	}

	nextCharacter, _ := reader.ReadByte()
	if nextCharacter != '\n' {
		return "", fmt.Errorf("request line not ending in LF")
	}

	trimmed := strings.TrimSuffix(maybeEndsInCR, "\r")
	return trimmed, nil
}

func (parser RFC7230RequestParser) OldParseRequest(reader *bufio.Reader) (*Request, *ParseError) {
	method, _ := readUpTo(reader, ' ')
	if len(method) == 0 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	} else if len(method) > maxLengthOfFieldInRequestLine {
		return nil, &ParseError{StatusCode: 501, Reason: "Not Implemented"}
	}

	target, _ := readUpTo(reader, ' ')
	if len(target) > maxLengthOfFieldInRequestLine {
		return nil, &ParseError{StatusCode: 414, Reason: "URI Too Long"}
	}

	version, _ := readUpTo(reader, '\r')
	return &Request{
		Method:  method,
		Target:  target,
		Version: version,
	}, nil
}

func readUpTo(reader *bufio.Reader, delimiter byte) (string, error) {
	field, err := reader.ReadString(delimiter)
	return strings.TrimSuffix(field, string(delimiter)), err
}

type Request struct {
	Method  string
	Target  string
	Version string
}

type ParseError struct {
	StatusCode int
	Reason     string
}
