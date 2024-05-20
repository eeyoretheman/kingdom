package Writers

import (
	"fmt"
	ConnClosed "kli/internal/conn_closed"
	Read "kli/internal/reader"
	"log"
	"net"
)

type Request struct {
	To      string
	Command string
	Body    string
}

type ReadWriter struct {
	Conn     net.Conn
	From     string
	Request  chan Request
	Response chan []byte
}

func ReadWriterHandler(writer *ReadWriter) {
	reader := make(chan []byte)
	readerErr := make(chan error)
	conn := writer.Conn
	request := writer.Request
	response := writer.Response

	go Read.Reader(conn, reader, readerErr)

	for {
		select {
		case input := <-request:
			err := error(nil)

			if input.To == "!" {
				_, err = conn.Write([]byte(fmt.Sprintf("%s ! %s %s", writer.From, input.Command, input.Body)))
			} else {
				_, err = conn.Write([]byte(fmt.Sprintf("%s %s %s %s", writer.From, input.To, input.Command, input.Body)))
			}

			if err != nil {
				log.Printf("Write failed; Error: %s\n", err)
			}
		case data := <-reader:
			response <- data
		case err := <-readerErr:
			log.Printf("Read failed; Error: %s\n", err)

			if ConnClosed.ConnClosed(err) {
				response <- []byte("!")
			}
		}
	}
}
