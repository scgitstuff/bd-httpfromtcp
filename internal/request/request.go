package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	all, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("ReadAll() failed")
	}

	s := string(all)
	lines := strings.Split(s, "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("RequestFromReader() bad request format")
	}

	x, err := parseRequestLine(lines[0])

	return x, err
}

func parseRequestLine(line string) (*Request, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad line format")
	}
	x := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	parts = strings.Split(x.HttpVersion, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("bad HttpVersion format")
	}
	if parts[0] != "HTTP" {
		return nil, fmt.Errorf("bad HttpVersion protocol, must be HTTP")
	}
	if parts[1] != "1.1" {
		return nil, fmt.Errorf("bad HttpVersion number, must be 1.1")
	}
	x.HttpVersion = parts[1]

	isMethodValid, err := regexp.Match("[A-Z]", []byte(x.Method))
	if err != nil {
		return nil, err
	}
	if !isMethodValid {
		return nil, fmt.Errorf("bad Method, must be upper case")
	}

	return &Request{x}, nil
}
