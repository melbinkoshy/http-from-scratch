package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{map[string]string{}}
}

var rn = []byte("\r\n")

var ErrInvalidHeader = fmt.Errorf("invalid header")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", errors.Join(fmt.Errorf("error inside parse header1"), ErrInvalidHeader)
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", errors.Join(fmt.Errorf("error inside parse header1"), ErrInvalidHeader)
	}

	name1 := strings.TrimLeft(string(name), " ")
	fmt.Println("The key for the header is ", name1)
	if !isValidKey(name1) {
		return "", "", errors.Join(fmt.Errorf("invalid characters"), ErrInvalidHeader)
	}

	return name1, string(value), nil
}

func isValidKey(a string) bool {
	allowedSpecialChar := "!#$%&'*+-.^_`|~"
	for _, s := range a {
		if unicode.IsLetter(s) || unicode.IsDigit(s) {
			continue
		}

		if strings.ContainsRune(allowedSpecialChar, s) {
			continue
		} else {
			return false
		}
	}
	return true
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}
		if idx == 0 {
			done = true
			//  read += len(rn)
			return read, done, nil
		}

		headerKey, headerValue, err := parseHeader(data[read : read+idx])
		read += idx + len(rn)

		if err != nil {
			return 0, done, err
		}
		h.Set(headerKey, headerValue)
	}

	return read, done, nil

}
