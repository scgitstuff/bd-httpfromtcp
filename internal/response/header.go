package response

import (
	"httpfromtcp/internal/headers"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Add("Content-Length", strconv.Itoa(contentLen))
	h.Add("Connection", "close")
	h.Add("Content-Type", "text/plain")

	return h
}
