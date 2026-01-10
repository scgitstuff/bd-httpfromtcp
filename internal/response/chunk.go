package response

import (
	"strconv"
)

func (w *Writer) WriteChunkedStart() error {
	h := GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Replace("Transfer-Encoding", "chunked")

	err := w.WriteStatusLine(StatusCodeSuccess)
	if err != nil {
		return err
	}

	return w.WriteHeaders(h)
}

func (w *Writer) WriteChunkedBody(chunk []byte) (int, error) {
	chunkLenHex := strconv.FormatInt(int64(len(chunk)), 16)
	// fmt.Printf("\n********************chunkLenHex: {%s}\n", chunkLenHex)

	n1, err := w.WriteBody([]byte(chunkLenHex + "\r\n"))
	if err != nil {
		return n1, err
	}

	line := string(chunk) + "\r\n"
	n2, err := w.WriteBody([]byte(line))

	return n1 + n2, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.WriteBody([]byte("0\r\n\r\n"))
}
