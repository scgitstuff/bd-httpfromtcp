package response

import "httpfromtcp/internal/headers"

type Writer struct {
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return 0, nil
}
