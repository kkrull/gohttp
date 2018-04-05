package http

import (
	"bufio"
	"strings"
	"fmt"
)

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
		Method:  fields[0],
		Target:  fields[1],
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

type Request struct {
	Method  string
	Target  string
	Version string
}

type ParseError struct {
	StatusCode int
	Reason     string
}
