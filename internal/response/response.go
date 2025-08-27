package response

import (
	"errors"
	"fmt"
	"io"
	"tcp_http/internal/headers"
)

type Response struct {
}

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case StatusOK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case StatusInternalError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		return errors.New("unknown status code")
	}

	_, err := w.writer.Write([]byte(statusLine))
	return err

}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{}

	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")
	_, err := w.writer.Write(b)
	return err

}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	return n, err

}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
