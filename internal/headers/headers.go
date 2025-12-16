package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	s := string(data)
	x := strings.Index(s, ":")
	if x == -1 {
		return 0, true, fmt.Errorf("':' char not found")
	}
	left := s[:x]
	if left[len(left)-1:] == " " {
		return 0, true, fmt.Errorf("field-name is invalid, space before ':'")
	}
	right := s[x+1:]
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	// fmt.Printf("\n****** '%s' : '%s'\n", left, right)

	// TODO: basic string parse done
	// add check for CRLF

	h[left] = right

	return 0, true, nil
}
