package response

import (
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int

const (
	good StatusCode = 200
	no   StatusCode = 400
	bad  StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	return nil
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	return nil
}
