package main

import (
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

			if i := bytes.IndexByte(data, '\n'); i == -1 {
				st += string(data[:n])
			} else {
				st += string(data[:i])
				ans <- st
				st = string(data[i+1 : n])
			}

		}

	}()

	return ans
}
func main() {

	listner, err := net.Listen("tcp", ":1764")
	if err != nil {
		log.Fatal("error :", err)
	}
	conn, err := listner.Accept()
	if err != nil {
		log.Fatal("error :", err)
	}

	lines := getLinesChannel(conn)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
