package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var CRLF string = "\r\n"

var Invalid_Key_Pair_Value = fmt.Errorf("Invalid key pair value")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(CRLF))
	// no crlf to start parsing header
	if crlfIdx == -1 {
		return 0, false, nil
	}
	// theres no header to start parsing it
	if crlfIdx == 0 {
		return 0, true, nil
	}

	colonIdx := bytes.Index(data, []byte(":"))

	if colonIdx == -1 {
		return 0, true, Invalid_Key_Pair_Value
	}
}

func formatData(data []byte) []byte {
	result := []byte
}
