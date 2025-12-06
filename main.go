package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const MSG_File = "messages.txt"

func main() {
	readLines(MSG_File)
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg)
	}
}

func readLines(fileName string) {
	fmt.Printf("Reading data from %s\n", fileName)
	fmt.Println("=====================================")

	f, err := os.Open(MSG_File)
	failOnErr(err, fmt.Sprintf("failed to open file: %s", MSG_File))
	defer f.Close()

	buf := make([]byte, 8)
	line := ""
	for {
		n, err := f.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					fmt.Printf("read: %s\n", line)
				}
				break
			}
			failOnErr(err, fmt.Sprintf("Read failed:\n%v\n", err))
		}

		line += string(buf[:n])
		lines := strings.Split(line, "\n")
		if len(lines) > 1 {
			fmt.Printf("read: %s\n", lines[0])
			line = lines[1]
			continue
		}
	}
}
