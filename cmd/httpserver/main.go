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

	// fmt.Println("**************************************************")
	// fmt.Println(req.String())
	// fmt.Println("**************************************************")
	out := server.HandlerError{}

	// TODO: I'll get back to finish this
	if req.RequestLine.RequestTarget == "/yourproblem" {
		out.StatusCode = response.NO
		out.Message = "400 Your problem is not my problem\n"
		return &out
	}

	// w.Write([]byte("STUFF"))

	return &server.HandlerError{StatusCode: response.BAD, Message: "This should not happen"}
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
