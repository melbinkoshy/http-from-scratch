package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type parserState string

const (
	stateInit  parserState = "init"
	stateDone  parserState = "done"
	stateError parserState = "errorState"
)

type Request struct {
	RequestLine RequestLine

	state parserState
}

func (r *Request) done() bool {
	return r.state == stateDone || r.state == stateError
}

func newRequest() *Request {
	return &Request{
		state: stateInit,
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) validHTTP() bool {
	return r.HttpVersion == "1.1"
}

var ErrBadReqLine = fmt.Errorf("invalid requestLine")
var ErrUnsupportedVersion = fmt.Errorf("upsupported http version")
var ErrRequestInErrState = fmt.Errorf("request in error state")

var SEPARATOR = []byte("\r\n")

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	//Note : Buffer could get overrun ... a header that exceeds 1k could do that
	//or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		//TODO : what to do with the errors
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}
	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case stateError:
			return 0, ErrRequestInErrState
		case stateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = stateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = stateDone
		case stateDone:
			break outer
		}

	}

	return read, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ErrBadReqLine
	}
	if string(parts[0]) != "GET" && string(parts[0]) != "POST" {
		return nil, 0, ErrBadReqLine
	}

	httpParts := strings.Split(string(parts[2]), "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" {
		return nil, 0, ErrBadReqLine
	}

	reqLine := &RequestLine{
		HttpVersion:   httpParts[1],
		RequestTarget: string(parts[1]),
		Method:        string(parts[0]),
	}

	if !reqLine.validHTTP() {
		return nil, 0, ErrUnsupportedVersion
	}

	return reqLine, read, nil
}
