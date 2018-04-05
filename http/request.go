package http

import (
	"bufio"
	"strings"
)

const maxLengthOfFieldInRequestLine = 8000

type RequestParser interface {
	ParseRequest(reader *bufio.Reader) (*Request, *ParseError)
}

type RFC7230RequestParser struct{}

func (parser RFC7230RequestParser) ParseRequest(reader *bufio.Reader) (*Request, *ParseError) {
	requestLineWithCR, _ := reader.ReadString('\r')
	requestLineWithCR = strings.TrimSuffix(requestLineWithCR, "\r")
	reader.ReadString('\n')

	fields := strings.Split(requestLineWithCR, " ")
	if len(fields) != 3 {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}

	method := fields[0]
	target := fields[1]
	version := fields[2]

	//End of headers
	reader.ReadString('\r')
	reader.ReadString('\n')

	return &Request{
		Method:  method,
		Target:  target,
		Version: version,
	}, nil
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
