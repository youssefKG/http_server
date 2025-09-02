package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

var CRLF string = "\r\n"
var KeySpecialChars string = "!#$%&'*+-^_`|~"

var Invalid_Key_Pair_Value = fmt.Errorf("Invalid key pair value")
var Invalid_Key_Format = fmt.Errorf("Invalid key format")
var Error_Missing_End_Headers = fmt.Errorf("Missing end of headers")
var Error_Malformed_Headers = fmt.Errorf("Malformed headers")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(CRLF))
	// no crlf to start parsing header
	if crlfIdx == -1 {
		return 0, false, nil
	}
	// theres no header to start parsing it
	if crlfIdx == 0 {
		return len(CRLF), true, nil
	}

	colonIdx := bytes.Index(data[:crlfIdx], []byte(":"))

	if colonIdx == -1 {
		return 0, false, Error_Malformed_Headers
	}

	formatedKeyBytes, err := formatKey(data[:colonIdx])
	if err != nil {
		return 0, false, err
	}
	key := string(formatedKeyBytes)
	value, err := formatValue(data[colonIdx+1 : crlfIdx])
	if err != nil {
		return 0, false, err
	}
	h.Set(key, string(value))
	return len(CRLF) + crlfIdx, false, nil
}

func formatValue(data []byte) ([]byte, error) {
	whiteSpacesAtStart := 0
	whitesSpacesAtEnd := 0
	dataLen := len(data)
	for whiteSpacesAtStart < dataLen &&
		unicode.IsSpace(rune(data[whiteSpacesAtStart])) {
		whiteSpacesAtStart++
	}
	for whitesSpacesAtEnd < dataLen-whiteSpacesAtStart &&
		unicode.IsSpace(rune(data[dataLen-1-whitesSpacesAtEnd])) {
		whitesSpacesAtEnd++
	}
	return data[whiteSpacesAtStart : len(data)-whitesSpacesAtEnd], nil
}

func formatKey(data []byte) ([]byte, error) {
	whiteSpacesAtStart := 0
	dataLen := len(data)
	for whiteSpacesAtStart < dataLen &&
		unicode.IsSpace(rune(data[whiteSpacesAtStart])) {
		whiteSpacesAtStart++
	}
	if whiteSpacesAtStart == dataLen {
		return nil, Invalid_Key_Format
	}

	key := data[whiteSpacesAtStart:]
	if len(key) < 1 {
		return nil, Invalid_Key_Format
	}
	for _, c := range key {
		if unicode.IsSpace(rune(c)) {
			return nil, Invalid_Key_Format
		}
	}
	if !isValidKey(key) {
		return nil, Invalid_Key_Format
	}
	return data[whiteSpacesAtStart:], nil
}

func NewHeaders() Headers {
	return Headers{}
}

func isValidKey(key []byte) bool {
	for _, c := range key {
		if !unicode.IsLetter(rune(c)) && !unicode.IsNumber(rune(c)) &&
			!strings.Contains(KeySpecialChars, string(c)) {
			return false
		}
	}
	return true
}

func (h Headers) Set(key string, value string) {
	lowerKey := strings.ToLower(key)
	oldValue, exist := h[lowerKey]
	if exist {
		h[lowerKey] = oldValue + ", " + value
		return
	}
	h[lowerKey] = value
}

func (h Headers) Get(key string) (string, bool) {
	value, exist := h[strings.ToLower(key)]
	return value, exist
}
