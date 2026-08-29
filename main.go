package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error :", err)
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
