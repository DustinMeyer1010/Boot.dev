package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/DustinMeyer1010/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42068")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, _ := listener.Accept()

		if conn == nil {
			continue
		}

		fmt.Println("Connection created ", conn.RemoteAddr())

		request, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Request Line:\n- Method: %s\n- Target: %s\n- Verison: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)

		fmt.Println("Headers: ")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)

		}

		fmt.Println("Body: ")
		fmt.Printf("- %s\n", string(request.Body))
		fmt.Println("Closed Connection to", conn.RemoteAddr())
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
