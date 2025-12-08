package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const PORT = "42069"
const SERVER = "127.0.0.1"

func main() {
	// TODO: clean up duplicate code, first pass just getting lesson OK
	_ = getLinesChannel

	runServer()
}

func connect() *net.UDPConn {
	address := fmt.Sprintf("%s:%s", SERVER, PORT)
	x, err := net.ResolveUDPAddr("udp", address)
	failOnErr(err, fmt.Sprintf("failed to resolve: %s", address))
	listener, err := net.DialUDP("udp", nil, x)
	failOnErr(err, fmt.Sprintf("failed to open: %s", address))
	fmt.Printf("Listening on: %s\n", listener.LocalAddr())

	return listener
}

func runServer() {

	listener := connect()
	defer listener.Close()

	buf := bufio.NewReader(os.Stdin)
	_ = buf

	for {
		fmt.Print("\n>")

		line, err := buf.ReadString('\n')
		failOnErr(err, "ReadString() failed")

		listener.Write([]byte(line))
		failOnErr(err, "Write() failed")
	}
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
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
