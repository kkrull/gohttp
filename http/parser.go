package http

import (
	"bufio"
	"strings"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

/* RequestMessageParser */

type RequestMessageParser struct {
	reader *bufio.Reader
}

func (parser *RequestMessageParser) Parse() (ok *RequestLine, err Response) {
	requestLine, err := parser.readCRLFLine()
	if err != nil {
		return nil, err
	}

	return parser.doParseRequestLine(requestLine)
}

func (parser *RequestMessageParser) doParseRequestLine(requestLine string) (ok *RequestLine, err Response) {
	requested, err := parser.parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}

	return parser.doParseHeaders(requested)
}

func (parser *RequestMessageParser) parseRequestLine(text string) (ok *RequestLine, badRequest Response) {
	fields := strings.Split(text, " ")
	if len(fields) != 3 {
		return nil, &clienterror.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	return &RequestLine{
		Method: fields[0],
		Target: fields[1],
	}, nil
}

func (parser *RequestMessageParser) doParseHeaders(requested *RequestLine) (ok *RequestLine, err Response) {
	err = parser.parseHeaderLines()
	if err != nil {
		return nil, err
	}

	return requested, nil
}

func (parser *RequestMessageParser) parseHeaderLines() (badRequest Response) {
	isBlankLineBetweenHeadersAndBody := func(line string) bool { return line == "" }

	for {
		line, err := parser.readCRLFLine()
		if err != nil {
			return err
		} else if isBlankLineBetweenHeadersAndBody(line) {
			return nil
		}
	}
}

func (parser *RequestMessageParser) readCRLFLine() (line string, badRequest Response) {
	maybeEndsInCR, _ := parser.reader.ReadString('\r')
	if len(maybeEndsInCR) == 0 {
		return "", &clienterror.BadRequest{DisplayText: "end of input before terminating CRLF"}
	} else if !strings.HasSuffix(maybeEndsInCR, "\r") {
		return "", &clienterror.BadRequest{DisplayText: "line in request header not ending in CRLF"}
	}

	nextCharacter, _ := parser.reader.ReadByte()
	if nextCharacter != '\n' {
		return "", &clienterror.BadRequest{DisplayText: "message header line does not end in LF"}
	}

	trimmed := strings.TrimSuffix(maybeEndsInCR, "\r")
	return trimmed, nil
}

type RequestLine struct {
	Method          string
	Target          string
	QueryParameters map[string]string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}
