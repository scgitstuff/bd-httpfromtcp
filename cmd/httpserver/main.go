package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, doStuff)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func doStuff(w io.Writer, req *request.Request) *server.HandlerError {
	out := server.HandlerError{}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		out.StatusCode = response.StatusCodeBadRequest
		out.Message = "Your problem is not my problem\n"
		return &out
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		out.StatusCode = response.StatusCodeInternalServerError
		out.Message = "Woopsie, my bad\n"
		return &out
	}

	w.Write([]byte("All good, frfr\n"))

	return nil
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
