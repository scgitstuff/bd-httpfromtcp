package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	GOOD StatusCode = 200
	NO   StatusCode = 400
	BAD  StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	phrases := map[StatusCode]string{
		GOOD: "HTTP/1.1 200 OK\r\n",
		NO:   "HTTP/1.1 400 Bad Request\r\n",
		BAD:  "HTTP/1.1 500 Internal Server Error\r\n",
	}

	v, ok := phrases[statusCode]
	if !ok {
		return fmt.Errorf("undefined status code: %d", statusCode)
	}

	// fmt.Print(v)

	_, err := w.Write([]byte(v))

	return err
}

func WriteBody(w io.Writer, body []byte) error {
	_, err := w.Write(body)
	// _, err = w.Write([]byte("\r\n"))

	return err
}
