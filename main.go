package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"unicode"
)

func getLinesChannel(c net.Conn) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer c.Close()
		var data strings.Builder
		buffer := make([]byte, 8)
		for {
			readedByte, err := c.Read(buffer)
			if err != nil {
				break
			}
			data.WriteString(string(buffer[:readedByte]))
		}
		stringParts := strings.Split(data.String(), "\n")
		for _, s := range stringParts {
			if s != "" {
				out <- s
			}
		}
	}()
	return out
}

// func main() {
// 	listener, err := net.Listen("tcp", ":42069")
//
// 	if err != nil {
// 		fmt.Print("connection error", err)
// 		return
// 	}
//
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Print("connect errro", err)
// 			return
// 		}
// 		lines := getLinesChannel(conn)
//
// 		for line := range lines {
// 			fmt.Printf("%s\n", line)
// 		}
// 	}
// }

// func formatValue(data []byte) ([]byte, error) {
// 	whiteSpacesAtStart := 0
// 	whitesSpacesAtEnd := 0
// 	dataLen := len(data)
// 	for whiteSpacesAtStart < dataLen &&
// 		unicode.IsSpace(rune(data[whiteSpacesAtStart])) {
// 		whiteSpacesAtStart++
// 	}
// 	for whitesSpacesAtEnd < dataLen-whiteSpacesAtStart &&
// 		unicode.IsSpace(rune(data[dataLen-1-whitesSpacesAtEnd])) {
// 		whitesSpacesAtEnd++
// 	}
// 	return data[whiteSpacesAtStart : len(data)-whitesSpacesAtEnd], nil
// }
//
// func main() {
// 	value, err := formatValue([]byte("localhost:42069\r\n\r\n"))
// 	if err != nil {
// 		fmt.Println("error occured")
// 		return
// 	}
// 	fmt.Println(string(value))
// }

var Invalid_Key_Format = fmt.Errorf("Invalid key format")

func formatKey(data []byte) ([]byte, error) {
	whiteSpacesAtStart := 0
	dataLen := len(data)
	for whiteSpacesAtStart < dataLen &&
		unicode.IsSpace(rune(data[whiteSpacesAtStart])) {
		whiteSpacesAtStart++
	}
	fmt.Println("hello", string(data))
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
	return data[whiteSpacesAtStart:], nil
}
func main() {
	header := []byte("Host: localhost:42069\r\n\r\n")
	colonIdx := bytes.Index(header, []byte(":"))
	if colonIdx == -1 {
		return
	}

	key, err := formatKey(header[:colonIdx])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(key))
}
