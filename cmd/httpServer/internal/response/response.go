package response

import (
	"fmt"
	"httpServer/cmd/httpServer/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

const CRLF string = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	httpStatusText := httpStatusText(statusCode)
	startLine := fmt.Sprintf("HTTP/1.1 %d %s %s", statusCode, httpStatusText, CRLF)
	_, err := w.Write([]byte(startLine))
	if err != nil {
		return err
	}
	return nil
}

func httpStatusText(StatusCode StatusCode) string {
	switch StatusCode {
	case Ok:
		return "ok"
	case BadRequest:
		return "Bad Request"
	case InternalServerError:
		return "Internal Server Error"
	}
	return ""
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		fieldLine := getFieldLine(key, value)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders["Content-Length"] = strconv.Itoa(contentLen)
	defaultHeaders["Content-Type"] = "text/plain"
	defaultHeaders["Connection"] = "close"
	return defaultHeaders
}

func getFieldLine(key string, value string) string {
	return fmt.Sprintf("%s: %s %s", key, value, CRLF)
}
