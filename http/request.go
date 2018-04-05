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
	//requestLine, _ := readCRLFLine(reader)
	//requestLineWithCR, _ := reader.ReadString('\r')
	//reader.ReadString('\n')
	//requestLine := strings.TrimSuffix(requestLineWithCR, "\r")

	//fields := strings.Split(requestLine, " ")
	//if len(fields) != 3 {
	//	return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	//}

	//method := fields[0]
	//target := fields[1]
	//version := fields[2]

	//End of headers
	_, err := readCRLFLine(reader)
	if err != nil {
		return nil, &ParseError{StatusCode: 400, Reason: "Bad Request"}
	}
	return nil, nil

	//return &Request{
	//	Method:  method,
	//	Target:  target,
	//	Version: version,
	//}, nil
}

func readCRLFLine(reader *bufio.Reader) (string, error) {
	requestLineWithCR, _ := reader.ReadString('\r')
	if len(requestLineWithCR) == 0 {
		return "", fmt.Errorf("request line not ending in CR")
	} else if requestLineWithCR[len(requestLineWithCR) -1:] != "\r" {
		return "", fmt.Errorf("request line not ending in CR")
	}

	shouldBeLF, _ := reader.ReadByte()
	if shouldBeLF != '\n' {
		return "", fmt.Errorf("request line not ending in LF")
	}

	requestLine := strings.TrimSuffix(requestLineWithCR, "\r")
	return requestLine, nil
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
