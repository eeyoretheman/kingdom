package main

import (
	// Consider termui instead
	"bufio"
	"fmt"
	writer "kli/internal/writer"
	"log"
	"net"
	"os"
	"strings"
)

func responseHandler(writer writer.Writer) {
	for {
		msg := <-writer.Response
		fmt.Print(string(msg))
	}
}

func main() {
	writers := make(map[string]*writer.Writer)
	current := ""

	serverConn, err := net.Dial("tcp", "localhost:2222")

	if err != nil {
		log.Printf("Could not connect to server; Error: %s\n", err)
		return
	}

	requestChan := make(chan writer.Request)
	responseChan := make(chan []byte)
	serverWriter := writer.Writer{Conn: serverConn, From: "!", Request: requestChan, Response: responseChan}

	id := ""

	go writer.WriterHandler(&serverWriter)

	id = strings.TrimSuffix(string(<-serverWriter.Response), "\n")
	serverWriter.From = id

	writers[id] = &serverWriter
	current = id

	stdinReader := bufio.NewReader(os.Stdin)

	go responseHandler(serverWriter)

	for {
		if current == "" {
			log.Printf("No current writer\n")
			return
		}

		line, _, _ := stdinReader.ReadLine()
		pair := strings.SplitN(string(line), " ", 2)

		if len(pair) < 2 {
			pair = append(pair, ".")
		}

		writers[current].Request <- writer.Request{To: "!", Command: pair[0], Body: pair[1]}
	}
}
