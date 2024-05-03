package clients

import (
	ConnClosed "kingdom/internal/conn_closed"
	Reader "kingdom/internal/reader"
	"log"
	"net"
	"strings"
)

type Request struct {
	From    string
	Target  string
	Command string
	Body    []byte
}

type Client struct {
	Input  chan []byte
	Output chan []byte
}

func ClientHandler(conn net.Conn, requests chan<- Request, client Client) {
	reader := make(chan []byte)
	readerErr := make(chan error)

	go Reader.Reader(conn, reader, readerErr)

	for {
		select {
		case input := <-client.Input:
			_, err := conn.Write(input)

			if err != nil {
				log.Printf("Write failed(maybe); Error: %s\n", err)
			}
		case data := <-reader:
			parts := strings.SplitN(string(data), " ", 4)

			if len(parts) < 4 {
				_, err := conn.Write([]byte("Malformed request.\nEach field is atleast one byte, space separated:\n[From] [Target] [Command] [Body]\n"))

				if err != nil {
					log.Printf("Write failed(maybe); Error: %s\n", err)
				}

				break
			}

			requests <- Request{From: parts[0], Target: parts[1], Command: parts[2], Body: []byte(parts[3])}
		case err := <-readerErr:
			log.Printf("Read failed; Error: %s\n", err)

			if ConnClosed.ConnClosed(err) {
				return
			}
		}
	}
}
