package main

import (
	"fmt"
	"net"
	"strings"
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

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Print("connection error", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("connect errro", err)
			return
		}
		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}
	}
}
