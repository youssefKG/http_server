package main

import (
	"fmt"
	"httpServer/cmd/httpServer/internal/request"
	"httpServer/cmd/httpServer/server"
	"io"
	"log"
	"os"
	"os/signal"
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

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	printRequest(req)
	w.Write([]byte("All good, frfr\n"))
	return nil
	return &server.HandlerError{
		StatusCode: 500,
		Message:    "Woopsie, my bad\n",
	}
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
