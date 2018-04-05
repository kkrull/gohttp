package http

import (
	"bufio"
	"strings"
)

const maxLengthOfFieldInRequestLine = 8000

type RequestParser interface {
	ParseRequest(reader *bufio.Reader) (Request, *ParseError)
}

type RFC7230RequestParser struct{}

func (parser RFC7230RequestParser) ParseRequest(reader *bufio.Reader) (Request, *ParseError) {
	method, _ := readFieldFromRequestLine(reader)
	if len(method) == 0 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	} else if len(method) > maxLengthOfFieldInRequestLine {
		return nil, &ParseError{StatusCode: 501, Reason: "Not Implemented"}
	}

	target, _ := readFieldFromRequestLine(reader)
	if len(target) > maxLengthOfFieldInRequestLine {
		return nil, &ParseError{StatusCode: 414, Reason: "URI Too Long"}
	}

	return nil, nil
}

func readFieldFromRequestLine(reader *bufio.Reader) (string, error) {
	field, err := reader.ReadString(' ')
	return strings.TrimSuffix(field, " "), err
}

type Request interface {
}

type ParseError struct {
	StatusCode int
	Reason     string
}
