package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

const SERVER = "127.0.0.1"

type Server struct {
	listener    net.Listener
	isClosed    atomic.Bool
	handlerFunc Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {

	address := fmt.Sprintf("%s:%d", SERVER, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Server listening on: %s\n", listener.Addr())

	server := Server{
		listener:    listener,
		handlerFunc: handlerFunc,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	if s.listener == nil {
		return nil
	}
	err := s.listener.Close()
	s.listener = nil

	return err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// fmt.Println("Server connection accepted")
		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	body := bytes.NewBuffer([]byte{})
	handlerErr := s.handlerFunc(body, r)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	b := body.Bytes()
	h := response.GetDefaultHeaders(len(b))
	writer := response.NewWriter(conn)
	writer.WriteStatusLine(response.StatusCodeSuccess)
	writer.WriteHeaders(h)
	writer.WriteBody(b)
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
