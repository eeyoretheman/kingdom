package tellers

import (
	ConnClosed "kingdom/internal/conn_closed"
	Reader "kingdom/internal/reader"
	"log"
	"net"
)

type Response struct {
	Target string
	Body   []byte
}

type Teller struct {
	Input  chan []byte
	Output chan []byte
	Owner  string
}

func TellerHandler(conn net.Conn, responses chan<- Response, teller *Teller) {
	reader := make(chan []byte)
	readerErr := make(chan error)

	go Reader.Reader(conn, reader, readerErr)

	for {
		select {
		case input := <-teller.Input:
			_, err := conn.Write(input)

			if err != nil {
				log.Printf("Write failed(maybe); Error: %s\n", err)
			}
		case data := <-reader:
			responses <- Response{Target: teller.Owner, Body: data}
		case err := <-readerErr:
			log.Printf("Read failed; Error: %s\n", err)

			if ConnClosed.ConnClosed(err) {
				responses <- Response{"!", []byte(teller.Owner)}
			}
		}
	}
}
