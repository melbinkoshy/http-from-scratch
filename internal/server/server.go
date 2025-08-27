package server

import (
	"fmt"
	"io"
	"net"
	"tcp_http/internal/request"
	"tcp_http/internal/response"
)

type Server struct {
	closed   bool
	handler  Handler
	listener net.Listener
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}
	s.handler(responseWriter, r)

}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, error := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if error != nil {
		return nil, error
	}
	server := &Server{
		closed:   false,
		handler:  handler,
		listener: listener,
	}
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	s.listener.Close()
	return nil
}

func (s *Server) Listen() {

}
