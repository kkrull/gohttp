package http

import (
	"bufio"
	"strings"

	"github.com/kkrull/gohttp/response/clientError"
	"github.com/kkrull/gohttp/response/serverError"
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
		return nil, err
	}

	fields := strings.Split(requestLine, " ")
	if len(fields) != 3 {
		return nil, &clientError.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	switch fields[0] {
	case "GET":
		return &GetRequest{
			BaseDirectory: parser.BaseDirectory,
			Target:        fields[1],
		}, nil
	default:
		return nil, &serverError.NotImplemented{Method: fields[0]}
	}
}

func parseHeaderLines(reader *bufio.Reader) Response {
	isBlankLineBetweenHeadersAndBody := func(line string) bool { return line == "" }

	for {
		line, err := readCRLFLine(reader)
		if err != nil {
			return err
		} else if isBlankLineBetweenHeadersAndBody(line) {
			return nil
		}
	}
}

func readCRLFLine(reader *bufio.Reader) (string, Response) {
	maybeEndsInCR, _ := reader.ReadString('\r')
	if len(maybeEndsInCR) == 0 {
		return "", &clientError.BadRequest{DisplayText: "end of input before terminating CRLF"}
	} else if !strings.HasSuffix(maybeEndsInCR, "\r") {
		return "", &clientError.BadRequest{DisplayText: "line in request header not ending in CRLF"}
	}

	nextCharacter, _ := reader.ReadByte()
	if nextCharacter != '\n' {
		return "", &clientError.BadRequest{DisplayText: "message header line does not end in LF"}
	}

	trimmed := strings.TrimSuffix(maybeEndsInCR, "\r")
	return trimmed, nil
}
