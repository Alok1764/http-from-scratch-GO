package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error :", err)
	}
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
			fmt.Printf("read: %s\n", st)
			st = string(data[i+1 : n])
		}

	}

}
