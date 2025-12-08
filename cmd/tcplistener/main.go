package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const MSG_File = "messages.txt"
const PORT = "42069"
const SERVER = "127.0.0.1"

func main() {
	// readLines(MSG_File)
	_ = readLines

	// f, err := os.Open(MSG_File)
	// failOnErr(err, fmt.Sprintf("failed to open file: %s", MSG_File))
	// defer f.Close()
	// ch := getLinesChannel(f)
	// for line := range ch {
	// 	fmt.Printf("read: %s\n", line)
	// }

	runServer()

}

func connect() net.Listener {
	address := fmt.Sprintf("%s:%s", SERVER, PORT)
	listener, err := net.Listen("tcp", address)
	failOnErr(err, fmt.Sprintf("failed to open: %s", address))
	fmt.Printf("Listening on: %s\n", listener.Addr())

	return listener
}

func runServer() {

	listener := connect()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		failOnErr(err, "Accept() failed")
		fmt.Println("connection accepted")

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}

		conn.Close()
		fmt.Println("connection closed")
	}
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
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

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)
		buf := make([]byte, 8)
		line := ""
		for {
			n, err := f.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					if line != "" {
						ch <- line
					}
					break
				}
				failOnErr(err, fmt.Sprintf("Read failed:\n%v\n", err))
			}
			line += string(buf[:n])

			lines := strings.Split(line, "\n")
			if len(lines) > 1 {
				for i := 0; i < len(lines)-1; i++ {
					ch <- lines[i]
				}
				line = lines[len(lines)-1]
			}
		}
	}()

	return ch
}
