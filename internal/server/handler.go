package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// TODO change sig, move handler code into httpserver
type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) String() string {
	return fmt.Sprintf("HandlerError {code: %d, msg: %s}", he.StatusCode, he.Message)
}

func (he HandlerError) Write(w io.Writer) {
	writer := response.NewWriter(w)
	writer.WriteStatusLine(he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	writer.WriteHeaders(headers)
	w.Write(messageBytes)
}
