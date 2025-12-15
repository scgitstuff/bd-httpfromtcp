package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const MSG_File = "messages.txt"
const PORT = "42069"
const SERVER = "127.0.0.1"

func main() {
	runServer()
}

func connect() net.Listener {
	address := fmt.Sprintf("%s:%s", SERVER, PORT)
	listener, err := net.Listen("tcp", address)
	failOnErr(err, fmt.Sprintf("failed to open: %s", address))
	fmt.Printf("Listening on: %s\n", listener.Addr())

	return listener
}

func runServer() {
	listener := connect()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		failOnErr(err, "Accept() failed")
		fmt.Println("connection accepted")

		r, err := request.RequestFromReader(conn)
		failOnErr(err, "RequestFromReader() failed")
		_ = r
		s := fmt.Sprintf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s",
			r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println(s)

		conn.Close()
		fmt.Println("connection closed")
	}
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
