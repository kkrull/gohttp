package http

import (
	"bufio"
	"strings"

	"github.com/kkrull/gohttp/msg/clienterror"
	"github.com/kkrull/gohttp/msg/servererror"
)

type RequestLineRouter struct {
	routes []Route
}

func (router *RequestLineRouter) AddRoute(route Route) {
	router.routes = append(router.routes, route)
}

func (router *RequestLineRouter) Routes() []Route {
	routes := make([]Route, len(router.routes))
	copy(routes, router.routes)
	return routes
}

func (router RequestLineRouter) ParseRequest(reader *bufio.Reader) (ok Request, routeError Response) {
	request, err := router.parseRequestLine(reader)
	if err != nil {
		return nil, err
	}

	headerError := parseHeaderLines(reader)
	if headerError != nil {
		return nil, headerError
	}

	return request, nil
}

func (router RequestLineRouter) parseRequestLine(reader *bufio.Reader) (Request, Response) {
	requestLineText, err := readCRLFLine(reader)
	if err != nil {
		return nil, err
	}

	requested, err := parseRequestLine(requestLineText)
	if err != nil {
		return nil, err
	}

	if request := router.routeRequest(requested); request != nil {
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

func (router RequestLineRouter) routeRequest(requested *RequestLine) Request {
	for _, route := range router.routes {
		request := route.Route(requested)
		if request != nil {
			return request
		}
	}

	return nil
}

type Route interface {
	Route(requested *RequestLine) Request
}

type RequestLine struct {
	Method      string
	Target      string
	QueryString string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}
