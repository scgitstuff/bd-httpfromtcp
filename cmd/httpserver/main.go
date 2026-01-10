package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, doStuff)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func doStuff(w *response.Writer, req *request.Request) {

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		respondCHUNK(w, strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/"))
		return
	}

	handleHTML(w, req)
}

func respondHTML(w *response.Writer, statusCode response.StatusCode, html string) {
	b := []byte(html)
	h := response.GetDefaultHeaders(len(b))
	h.Replace("Content-Type", "text/html")
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
	w.WriteBody(b)
}

func respondTEXT(w *response.Writer, statusCode response.StatusCode, body string) {
	b := []byte(body)
	h := response.GetDefaultHeaders(len(b))
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
	w.WriteBody(b)
}

func respondCHUNK(w *response.Writer, path string) {
	base := "https://httpbin.org"
	url := base + "/" + path
	const CHUNK = 1024
	buff := make([]byte, CHUNK)

	// fmt.Println(url)

	resp, err := http.Get(url)

	err = w.WriteChunkedStart()
	if err != nil {
		fmt.Printf("WriteChunkedStart: bad stuff happened:\n\t%v\n", err)
	}

	for {
		n, err := resp.Body.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("respondCHUNK: read error:\n\t%v\n", err)
			break
		}
		// fmt.Printf("respondCHUNK: read: %d\n", n)

		if n <= 0 {
			break
		}

		x, err := w.WriteChunkedBody(buff[:n])
		if err != nil {
			fmt.Printf("WriteChunkedBody: error:\n\t%v\n", err)
		}
		_ = x
		// fmt.Printf("respondCHUNK: write: %d\n", x)
	}

	w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("WriteChunkedBodyDone: bad stuff happened:\n\t%v\n", err)
	}
}

func failOnErr(err error, msg string) {
	if err != nil {
		// panic(err)
		log.Fatal(msg, "\n\t", err)
	}
}
