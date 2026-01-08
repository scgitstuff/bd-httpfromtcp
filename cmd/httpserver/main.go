package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, doStuffHTML)
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

func respondHTML(w *response.Writer, statusCode response.StatusCode, html string) {
	b := []byte(html)
	h := response.GetDefaultHeaders(len(b))
	h.Replace("Content-Type", "text/html")
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
	w.WriteBody(b)
}

func respondTEXT(w *response.Writer, statusCode response.StatusCode, body string) {
	b := []byte(body)
	h := response.GetDefaultHeaders(len(b))
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
	w.WriteBody(b)
}

func doStuff(w *response.Writer, req *request.Request) {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		respondTEXT(w, response.StatusCodeBadRequest,
			"Your problem is not my problem\n")
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		respondTEXT(w, response.StatusCodeInternalServerError,
			"Woopsie, my bad\n")
		return
	}

	respondTEXT(w, response.StatusCodeSuccess, "All good, frfr\n")
}

func doStuffHTML(w *response.Writer, req *request.Request) {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		s := `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
		respondHTML(w, response.StatusCodeBadRequest, s)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		s := `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
		respondHTML(w, response.StatusCodeInternalServerError, s)
		return
	}

	s := `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
	respondHTML(w, response.StatusCodeSuccess, s)
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
