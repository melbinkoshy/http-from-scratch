package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"tcp_http/internal/headers"
	"tcp_http/internal/request"
	"tcp_http/internal/response"
	"tcp_http/internal/server"
)

const port = 42069

func toStr(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%x", b)
	}
	return out
}

func Respond400() []byte {
	return []byte(`
		<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>

	`)
}

func Respond500() []byte {
	return []byte(`
		<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>

	`)
}

func Respond200() []byte {
	return []byte(`
		<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>

	`)
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := Respond200()
		status := response.StatusOK
		if req.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusBadRequest
			body = Respond400()

		} else if req.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInternalError
			body = Respond500()

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				body = Respond500()
				status = response.StatusInternalError
				w.WriteStatusLine(status)
				h.Replace("content-length", fmt.Sprintf("%d", len(body)))
				h.Replace("Content-Type", "text/html")
				w.WriteHeaders(*h)
				w.WriteBody(body)
				return // <--- THIS return is critical
			} else {
				w.WriteStatusLine(response.StatusOK)
				h.Delete("content-length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				w.WriteHeaders(*h)
				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break

					}
					fullBody = append(fullBody, data[:n]...)
					w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailer := headers.NewHeaders()
				out := sha256.Sum256(fullBody)
				trailer.Set("X-Content-SHA256", toStr(out[:]))
				trailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				w.WriteHeaders(*trailer)
				return
			}
		} else if req.RequestLine.RequestTarget == "/video" {
			file, err := os.ReadFile("./assets/nature.mp4")
			if err != nil {
				body = Respond500()
				status = response.StatusInternalError
				w.WriteStatusLine(status)
				h.Replace("content-length", fmt.Sprintf("%d", len(body)))
				h.Replace("Content-Type", "text/html")
				w.WriteHeaders(*h)
				w.WriteBody(body)
				return
			}
			h.Replace("Content-Type", "video/mp4")
			h.Replace("content-length", fmt.Sprintf("%d", len(file)))
			w.WriteStatusLine(response.StatusOK)
			w.WriteHeaders(*h)
			w.WriteBody(file)
			return
		}
		h.Replace("content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(body)

	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
