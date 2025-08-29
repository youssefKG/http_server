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
		fmt.Println("request version: ", r.RequestLine.HttpVersion)
		fmt.Println("request method: ", r.RequestLine.Method)
		fmt.Println("request target: ", r.RequestLine.RequestTarget)

	}
}
