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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) String() string {
	return fmt.Sprintf("HandlerError {code: %d, msg: %s}", he.StatusCode, he.Message)
}
