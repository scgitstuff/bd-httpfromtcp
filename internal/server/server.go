package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

const SERVER = "127.0.0.1"

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	address := fmt.Sprintf("%s:%d", SERVER, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Server listening on: %s\n", listener.Addr())

	server := Server{
		listener: listener,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
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
			if s.closed.Load() {
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
	out := `HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!`

	conn.Write([]byte(out))
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
