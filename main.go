package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const MSG_File = "messages.txt"

func main() {
	readMessages(MSG_File)
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg)
	}
}

func readMessages(fileName string) {
	fmt.Printf("Reading data from %s\n", fileName)
	fmt.Println("=====================================")

	f, err := os.Open(MSG_File)
	failOnErr(err, fmt.Sprintf("failed to open file: %s", MSG_File))
	defer f.Close()

	buf := make([]byte, 8)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			failOnErr(err, fmt.Sprintf("Read failed:\n%v\n", err))
		}
		fmt.Printf("read: %s\n", buf[:n])
	}
}
