package http

import (
	"bufio"
	"strings"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

type RFC7230RequestParser struct {
	Routes []Route
}

func (parser *RFC7230RequestParser) AddRoute(route Route) {
	parser.Routes = append(parser.Routes, route)
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

func (parser RFC7230RequestParser) parseRequestLine(reader *bufio.Reader) (Request, Response) {
	requestLineText, err := readCRLFLine(reader)
	if err != nil {
		return nil, err
	}

	requested, err := parseRequestLine(requestLineText)
	if err != nil {
		return nil, err
	}

	if request := parser.routeRequest(requested); request != nil {
		return request, nil
	}

	return nil, requested.NotImplemented()
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
		return "", &clienterror.BadRequest{DisplayText: "end of input before terminating CRLF"}
	} else if !strings.HasSuffix(maybeEndsInCR, "\r") {
		return "", &clienterror.BadRequest{DisplayText: "line in request header not ending in CRLF"}
	}

	nextCharacter, _ := reader.ReadByte()
	if nextCharacter != '\n' {
		return "", &clienterror.BadRequest{DisplayText: "message header line does not end in LF"}
	}

	trimmed := strings.TrimSuffix(maybeEndsInCR, "\r")
	return trimmed, nil
}

func parseRequestLine(text string) (*RequestLine, Response) {
	fields := strings.Split(text, " ")
	if len(fields) != 3 {
		return nil, &clienterror.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	return &RequestLine{
		Method: fields[0],
		Target: fields[1],
	}, nil
}

func (parser RFC7230RequestParser) routeRequest(requested *RequestLine) Request {
	for _, route := range parser.Routes {
		request := route.Route(requested.Method, requested.Target)
		if request != nil {
			return request
		}
	}

	return nil
}

type RequestLine struct {
	Method string
	Target string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}

type Route interface {
	Route(method string, target string) Request
}
