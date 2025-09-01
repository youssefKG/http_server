package server

import (
	"fmt"
	"httpServer/cmd/httpServer/internal/request"
	"httpServer/cmd/httpServer/internal/response"
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

type Handler func(w *response.Writer, req *request.Request)

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
	newRequest, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	writer := response.NewWriter()
	s.handler(writer, newRequest)
	conn.Write(writer.GetBuffer())
}

func NewServer(listener net.Listener, handler Handler) *Server {
	return &Server{
		listener: listener,
		handler:  handler,
	}
}
