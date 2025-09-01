package server

import (
	"bytes"
	"fmt"
	"httpServer/cmd/httpServer/internal/request"
	"httpServer/cmd/httpServer/internal/response"
	"io"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	colsed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := NewServer(listener, handler)
	if err != nil {
		return nil, err
	}
	go server.listen()
	return server, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.colsed.Load() {
				return
			}
			fmt.Println("accept error: ", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.colsed.Store(false)
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	r, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	var buffer bytes.Buffer
	var handlerBuffer bytes.Buffer
	handlerError := s.handler(&handlerBuffer, r)
	if handlerError != nil {
		response.WriteStatusLine(&buffer, handlerError.StatusCode)
		buffer.Write(fmt.Appendf(nil, "message: %s%s%s", handlerError.Message,
			response.CRLF, response.CRLF))
	} else {
		defaultHeaders := response.GetDefaultHeaders(handlerBuffer.Len())
		response.WriteStatusLine(&buffer, 200)
		response.WriteHeaders(&buffer, defaultHeaders)
		buffer.Write(handlerBuffer.Bytes())
	}
	conn.Write(buffer.Bytes())
}

func NewServer(listener net.Listener, handler Handler) *Server {
	return &Server{
		listener: listener,
		handler:  handler,
	}
}
