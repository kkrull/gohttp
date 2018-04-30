package http

import (
	"bufio"
	"strings"

	"github.com/kkrull/gohttp/msg/clienterror"
)

// Parses an HTTP request message one line at a time.
type LineRequestParser struct{}

func (parser *LineRequestParser) Parse(reader *bufio.Reader) (ok *requestMessage, err Response) {
	methodObject := &parseMethodObject{reader: reader}
	return methodObject.Parse()
}

type parseMethodObject struct {
	reader *bufio.Reader
}

func (parser *parseMethodObject) Parse() (ok *requestMessage, err Response) {
	requestLine, err := parser.readCRLFLine()
	if err != nil {
		return nil, err
	}

	return parser.doParseRequestLine(requestLine)
}

func (parser *parseMethodObject) doParseRequestLine(requestLine string) (ok *requestMessage, err Response) {
	requested, err := parser.parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}

	return parser.doParseHeaders(requested)
}

func (parser *parseMethodObject) doParseHeaders(requested *requestMessage) (ok *requestMessage, err Response) {
	err = parser.parseHeaders()
	if err != nil {
		return nil, err
	}

	return requested, nil
}

func (parser *parseMethodObject) parseRequestLine(text string) (ok *requestMessage, badRequest Response) {
	fields := strings.Split(text, " ")
	if len(fields) != 3 {
		return nil, &clienterror.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	return &requestMessage{
		method: fields[0],
		target: fields[1],
	}, nil
}

func (parser *parseMethodObject) parseHeaders() (badRequest Response) {
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

func (parser *parseMethodObject) readCRLFLine() (line string, badRequest Response) {
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
