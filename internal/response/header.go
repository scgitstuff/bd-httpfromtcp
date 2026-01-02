package response

import (
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	w.Write([]byte(headers.String()))
	_, err := w.Write([]byte("\r\n"))

	return err
}
