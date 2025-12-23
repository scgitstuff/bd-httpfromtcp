package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	state   requestState
	bodyLen int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const CR_LF = "\r\n"
const BUFF_SIZE = 8

func (rl *RequestLine) String() string {
	s := fmt.Sprintf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
		rl.Method, rl.RequestTarget, rl.HttpVersion)
	return s
}

func (r *Request) String() string {
	var stuff strings.Builder
	stuff.WriteString("Headers:\n")

	for k, v := range r.Headers {
		fmt.Fprintf(&stuff, "- %s: %s\n", k, v)
	}

	fmt.Fprintf(&stuff, "Body:\n%s\n", r.Body)

	return r.RequestLine.String() + stuff.String()
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFF_SIZE)
	readToIndex := 0
	req := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.state, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(CR_LF))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}
	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// Need more data
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// Need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, isDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if isDone {
			r.state = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		n := len(data)
		head, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.state = requestStateDone
			return n, nil
		}
		contentLen, err := strconv.Atoi(head)
		if err != nil {
			return 0, fmt.Errorf("bad 'Content-Length' header: %s", head)
		}
		// not speced; I don't consider it an error, just a useless header
		// if contentLen <= 0 {
		// 	r.state = requestStateDone
		// 	return 0, nil
		// }

		r.Body = append(r.Body, data...)
		r.bodyLen = len(r.Body)
		if r.bodyLen > contentLen {
			return 0, fmt.Errorf("Body length %d is longer than 'Content-Length' %d",
				r.bodyLen, contentLen,
			)
		}
		if r.bodyLen == contentLen {
			r.state = requestStateDone
			fmt.Println("I eated all the Body data")
			return n, nil
		}
		// TODO: why does this happen twice for each chunk?
		// fmt.Printf("************************%s\n", r.Body)
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
