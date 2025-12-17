package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const CR_LF = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(CR_LF))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	s := string(data[:idx])
	n = len(s) + 2
	key, value, found := strings.Cut(s, ":")
	if !found {
		return 0, false, fmt.Errorf("':' char not found")
	}

	if key[len(key)-1:] == " " {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	h[key] = value

	return n, false, nil
}
