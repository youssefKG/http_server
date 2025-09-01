package main

import (
	"fmt"
	"httpServer/cmd/tcplistener/internal/request"
	"net"
)

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
		r, err := request.RequestFromReader(conn)
		if err != nil {
			break
		}

		fmt.Println("Request line : ")
		fmt.Println("- Version: ", r.RequestLine.HttpVersion)
		fmt.Println("- Method: ", r.RequestLine.Method)
		fmt.Println("- Target: ", r.RequestLine.RequestTarget)
		fmt.Println("Headers: ")
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body: ")
		fmt.Println(string(r.Body))
	}
}
