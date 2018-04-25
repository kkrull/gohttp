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

func (router RequestLineRouter) ParseRequest(reader *bufio.Reader) (ok Request, err Response) {
	requestLine, err := readCRLFLine(reader)
	if err != nil {
		return nil, err
	}

	return router.doParseRequestLine(reader, requestLine)
}

func (router *RequestLineRouter) doParseRequestLine(reader *bufio.Reader, requestLine string) (ok Request, err Response) {
	requested, err := parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}

	return router.doParseHeaders(reader, requested)
}

func parseRequestLine(text string) (ok *RequestLine, badRequest Response) {
	fields := strings.Split(text, " ")
	if len(fields) != 3 {
		return nil, &clienterror.BadRequest{DisplayText: "incorrectly formatted or missing request-line"}
	}

	return &RequestLine{
		Method: fields[0],
		Target: fields[1],
	}, nil
}

func (router *RequestLineRouter) doParseHeaders(reader *bufio.Reader, requested *RequestLine) (ok Request, err Response) {
	err = parseHeaderLines(reader)
	if err != nil {
		return nil, err
	}

	return router.routeRequest(requested)
}

func parseHeaderLines(reader *bufio.Reader) (badRequest Response) {
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

func readCRLFLine(reader *bufio.Reader) (line string, badRequest Response) {
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

func (router RequestLineRouter) routeRequest(requested *RequestLine) (ok Request, notImplemented Response) {
	for _, route := range router.routes {
		request := route.Route(requested)
		if request != nil {
			return request, nil
		}
	}

	return nil, requested.NotImplemented()
}

type Route interface {
	Route(requested *RequestLine) Request
}

type RequestLine struct {
	Method          string
	Target          string
	QueryParameters map[string]string
}

func (requestLine *RequestLine) NotImplemented() Response {
	return &servererror.NotImplemented{Method: requestLine.Method}
}
