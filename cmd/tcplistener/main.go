package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"tcp_http/internal/request"
)

func main() {

	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		fmt.Println("New connection from:", conn.RemoteAddr())
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Request line")
		fmt.Println("-Method:", request.RequestLine.Method)
		fmt.Println("-Target:", request.RequestLine.RequestTarget)
		fmt.Println("-Version:", request.RequestLine.HttpVersion)
		fmt.Println("Headers")
		request.Headers.ForEach(func(n, v string) {
			fmt.Printf(" -%s: %s\n", n, v)
		})
		fmt.Println("Body")
		fmt.Printf("%s\n", request.Body)
	}

}

func readFromFile() {
	file, error := os.Open("./messages.txt")
	if error != nil {
		fmt.Println("Error opening file:", error)
		return
	}
	defer file.Close()

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
