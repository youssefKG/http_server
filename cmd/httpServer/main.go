package main

import (
	"fmt"
	"httpServer/cmd/httpServer/internal/headers"
	"httpServer/cmd/httpServer/internal/request"
	"httpServer/cmd/httpServer/internal/response"
	"httpServer/cmd/httpServer/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const port = 42069

func printRequest(r *request.Request) {
	fmt.Printf(
		"%s %s %s\n",
		r.RequestLine.Method,
		r.RequestLine.RequestTarget,
		r.RequestLine.HttpVersion)
	for key, value := range r.Headers {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println()
	fmt.Println(string(r.Body))
	fmt.Println()
}

func handler(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(200)
	newHeaders := headers.NewHeaders()
	// newHeaders["Content-Type"] = "text/html"
	fmt.Println("Method: ", req.RequestLine.Method)
	fmt.Println("Method: ", req.RequestLine.RequestTarget)
	newHeaders["Content-Type"] = "application/json"
	data, err := os.ReadFile("./test.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	newHeaders["Content-Length"] = strconv.Itoa(len(data))
	err = w.WriteHeaders(newHeaders)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteBody(data)
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
