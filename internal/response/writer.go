package response

import (
	"httpfromtcp/internal/headers"
	"io"
)

type Writer struct {
	ioW io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{ioW: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := w.ioW.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	w.ioW.Write([]byte(headers.String()))
	_, err := w.ioW.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	n, err := w.ioW.Write(body)
	return n, err
}
