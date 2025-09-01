package response

import (
	"bytes"
	"fmt"
	"httpServer/cmd/httpServer/internal/headers"
	"strconv"
)

type StatusCode int
type WriterState string

const (
	InitState              WriterState = "init"
	StatusLineStateWritten WriterState = "status_line"
	HeadersStateWritten    WriterState = "headers"
	BodyStateWritten       WriterState = "body"
)
const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

type Writer struct {
	state  WriterState
	buffer bytes.Buffer
}

var ERROR_STATUS_LINE_STATE = fmt.Errorf("Status code already written or wrong order")
var ERROR_BODY_STATE = fmt.Errorf("Must write headers before body")
var ERROR_HEADERS_STATE = fmt.Errorf("Must write status line before headers")

const CRLF string = "\r\n"

func (writer *Writer) WriteStatusLine(statusCode StatusCode) error {
	if writer.state != InitState {
		return ERROR_STATUS_LINE_STATE
	}
	httpStatusText := httpStatusText(statusCode)
	startLine := fmt.Sprintf("HTTP/1.1 %d %s %s", statusCode, httpStatusText, CRLF)
	_, err := writer.buffer.Write([]byte(startLine))
	if err != nil {
		return err
	}
	writer.state = StatusLineStateWritten
	return nil
}

func (writer *Writer) WriteHeaders(headers headers.Headers) error {
	if writer.state != StatusLineStateWritten {
		return ERROR_HEADERS_STATE
	}
	for key, value := range headers {
		fieldLine := getFieldLine(key, value)
		_, err := writer.buffer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	_, err := writer.buffer.Write([]byte(CRLF))
	if err != nil {
		return err
	}
	writer.state = HeadersStateWritten
	return nil
}

func (writer *Writer) WriteBody(body []byte) error {
	if writer.state != HeadersStateWritten {
		return ERROR_BODY_STATE
	}
	_, err := writer.buffer.Write(body)
	if err != nil {
		return err
	}
	writer.state = BodyStateWritten
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
	return fmt.Sprintf("%s: %s%s", key, value, CRLF)
}

func NewWriter() *Writer {
	return &Writer{
		state: InitState,
	}
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

func (writer *Writer) GetBuffer() []byte {
	return writer.buffer.Bytes()
}
