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
	return methodObject.ReadingRequestLine()
}

//A state machine that parses an HTTP request during the process of reading the request from input
type parseMethodObject struct {
	reader *bufio.Reader
}

func (parser *parseMethodObject) ReadingRequestLine() (ok *requestMessage, badRequest Response) {
	requestLine, err := parser.readCRLFLine()
	if err != nil {
		return nil, err
	}

	return parser.parsingRequestLine(requestLine)
}

func (parser *parseMethodObject) parsingRequestLine(requestLine string) (ok *requestMessage, badRequest Response) {
	fields := strings.Split(requestLine, " ")
	if len(fields) != 3 {
		return nil, &clienterror.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	return parser.parsingTarget(fields[0], fields[1])
}

func (parser *parseMethodObject) parsingTarget(method, target string) (ok *requestMessage, badRequest Response) {
	path, query, _ := splitTarget(target)
	requested := &requestMessage{
		method: method,
		target: target,
		path:   path,
	}

	return parser.parsingQueryString(requested, query)
}

func (parser *parseMethodObject) parsingQueryString(requested *requestMessage, rawQuery string) (ok *requestMessage, badRequest Response) {
	if len(rawQuery) == 0 {
		return parser.readingHeaders(requested)
	}

	stringParameters := strings.Split(rawQuery, "&")
	for _, stringParameter := range stringParameters {
		nameValueFields := strings.Split(stringParameter, "=")
		if len(nameValueFields) == 1 {
			requested.AddQueryFlag(nameValueFields[0])
		} else {
			requested.AddQueryParameter(nameValueFields[0], nameValueFields[1])
		}
	}

	return parser.readingHeaders(requested)
}

func (parser *parseMethodObject) readingHeaders(requested *requestMessage) (ok *requestMessage, badRequest Response) {
	isBlankLineBetweenHeadersAndBody := func(line string) bool { return line == "" }

	for {
		line, err := parser.readCRLFLine()
		if err != nil {
			return nil, err
		} else if isBlankLineBetweenHeadersAndBody(line) {
			return requested, nil
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

func splitTarget(target string) (path, query, fragment string) {
	splitOnQuery := strings.Split(target, "?")
	if len(splitOnQuery) == 1 {
		query = ""
		path, fragment = extractFragment(splitOnQuery[0])
		return
	}

	path = splitOnQuery[0]
	query, fragment = extractFragment(splitOnQuery[1])
	return
}

func extractFragment(target string) (prefix string, fragment string) {
	fields := strings.Split(target, "#")
	if len(fields) == 1 {
		return fields[0], ""
	} else {
		return fields[0], fields[1]
	}
}
