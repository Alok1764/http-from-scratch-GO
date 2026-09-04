package main

import (
	"alok/internal/request"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ans := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(ans)

		st := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i == -1 {
				st += string(data[:n])
			} else {
				st += string(data[:i])
				ans <- st
				st = string(data[i+1 : n])
			}

		}
		if st != "" {
			ans <- st
		}

	}()

	return ans
}
func main() {

	// on server go run ./cmd/tcplistener | tee /tmp/rawpost.http
	// on client curl.exe -X POST -H "Content-Type: application/json" -d '{"flavor":"dark mode"}' http://localhost:1764/coffee

	listner, err := net.Listen("tcp", ":1764")
	if err != nil {
		log.Fatal("error :", err)
	}
	conn, err := listner.Accept()
	if err != nil {
		log.Fatal("error :", err)
	}

	RequestLine, err := request.RequestFromReader(conn)

	fmt.Println("- Method: ", RequestLine.RequestLine.Method)

	fmt.Println("- Target: ", RequestLine.RequestLine.RequestTarget)

	fmt.Println("- version: ", RequestLine.RequestLine.HttpVersion)

}
