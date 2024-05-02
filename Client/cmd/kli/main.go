package main

import (
	// Consider termui instead
	writer "kli/internal/Writers"
	reader "kli/internal/reader"
	"log"
	"net"
)

func main() {
	writers := make(map[string]*writer.Writer)
	current := ""

	temp_conn, err := net.Dial("tcp", "localhost:2222")
	if err != nil {
		log.Printf("Could not connect to server; Error: %s\n", err)
		return
	}
	temp_request := make(chan writer.Request)
	temp_response := make(chan []byte)
	temp_writer := writer.Writer{Conn: temp_conn, From: "!", Request: temp_request, Response: temp_response}
	id := ""
	read := make(chan []byte)
	readerErr := make(chan error)
	reader.Reader(temp_conn, read, readerErr)
	id = string(<-read)
	temp_writer.From = id
	writers[id] = &temp_writer
	current = id
	go writer.WriterHandler(&temp_writer)

	for {
		if current == "" {
			log.Printf("No current writer\n")
			return
		}
	}
}
