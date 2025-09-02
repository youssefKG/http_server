package request

import (
	"fmt"
	"httpServer/cmd/httpServer/internal/headers"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
var ERROR_PARTIAL_CONTENT = fmt.Errorf("Partial Content")

var SEPARATOR string = "\r\n"

// statusLine: Method(Get, Post) path(/login) http/version(1.1) \r\n
// headers: Content-Type: text/html \r\n
// cookies:  token, token
// Content-Length: 123
// \n\r
// body

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
		if consumedBytes == 0 {
			continue
		} else {
			copy(buffer, buffer[consumedBytes:readtToIndex])
			readtToIndex -= consumedBytes
		}
	}

	bodySize, _ := strconv.Atoi(request.Headers["content-length"])
	if len(request.Body) != bodySize {
		return nil, ERROR_PARTIAL_CONTENT
	}
	return request, nil
}

func (r *Request) parseSingleLine(data []byte) (int, error) {
	switch r.state {
	case StateInit:
		requestLine, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
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
			r.state = StateBody
		}
		return n, nil
	case StateBody:
		_, exist := r.Headers.Get("content-length")
		if !exist {
			r.state = StateDone
			return 0, nil
		}
		r.Body = append(r.Body, data...)
		return len(data), nil
	default:
		return 0, fmt.Errorf("dont know to do yet")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != StateDone {
		n, err := r.parseSingleLine(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
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
