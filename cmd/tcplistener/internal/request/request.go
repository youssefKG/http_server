package request

import (
	"fmt"
	"httpServer/cmd/tcplistener/internal/headers"
	"io"
	"strings"
	"unicode"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	body        []byte
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_INVAlID_REQUEST_LINE = fmt.Errorf("invalid request line format.")
var ERROR_REQUEST_METHOD = fmt.Errorf("Excepect request method contain only capital lettres.")
var ERROR_HTTP_VERSION_FORMAT = fmt.Errorf("We only Support HTTP/1.1 for now.")
var ERROR_HTTP_VERSION = fmt.Errorf("Invalid http version")
var ERROR_INVALID_REQUEST_METHOD = fmt.Errorf("Invalid request method")
var Error_MISSING_END_HEADERS = fmt.Errorf("Missing End Error")

var SEPARATOR string = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	readtToIndex := 0
	buffer := make([]byte, 1024)
	request.Headers = headers.NewHeaders()

	for request.state != StateDone {
		numsBytesRead, err := reader.Read(buffer[readtToIndex:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		readtToIndex += numsBytesRead
		consumedBytes, err := request.parse(buffer[:readtToIndex])
		if err != nil {
			return nil, err
		}
		if consumedBytes != 0 {
			copy(buffer, buffer[:readtToIndex])
			readtToIndex -= consumedBytes
		}

	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case StateInit:
		requestLine, parsedBytes, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if parsedBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = StateHeaders
		return parsedBytes, nil
	case StateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done == true {
			r.state = StateDone
			return n, nil
		}
		return n, nil
	default:
		return 0, nil
	}
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func parseRequestLine(request string) (*RequestLine, int, error) {
	idx := strings.Index(request, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	if idx == -1 {
		return nil, 0, ERROR_INVAlID_REQUEST_LINE
	}
	requestLineParts := strings.Split(request[:idx], " ")
	if len(requestLineParts) != 3 {
		return nil, 0, ERROR_INVAlID_REQUEST_LINE
	}
	method := requestLineParts[0]
	for _, c := range method {
		if unicode.IsLower(c) {
			return nil, 0, ERROR_REQUEST_METHOD
		}
	}
	if method != "GET" && method != "POST" {
		return nil, 0, ERROR_REQUEST_METHOD
	}
	target := requestLineParts[1]
	httpVersion := requestLineParts[2]
	idx2 := strings.Index(httpVersion, "/")
	if idx2 == -1 || httpVersion[:idx2] != "HTTP" {
		return nil, 0, ERROR_HTTP_VERSION_FORMAT
	}
	if httpVersion[idx2+1:] != "1.1" {
		return nil, 0, ERROR_HTTP_VERSION
	}
	return &RequestLine{
			Method:        method,
			RequestTarget: target,
			HttpVersion:   httpVersion[idx2+1:],
		},
		idx + len(SEPARATOR),
		nil
}

func (r *Request) parseSingleLine(data []byte) (int, error) {
	switch r.state {
	case StateInit:
		requestLine, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		r.RequestLine = *requestLine
		r.state = StateHeaders
		return n, nil
	case StateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = StateDone
		}
		return n, nil
	default:
		return 0, fmt.Errorf("dont know to do yet")
	}
}
