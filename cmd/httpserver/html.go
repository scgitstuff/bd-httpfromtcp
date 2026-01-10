package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

func handleHTML(w *response.Writer, req *request.Request) {

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
