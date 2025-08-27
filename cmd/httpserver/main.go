package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tcp_http/internal/request"
	"tcp_http/internal/response"
	"tcp_http/internal/server"
)

const port = 42069

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

		}
		w.WriteStatusLine(status)
		h.Replace("content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
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
