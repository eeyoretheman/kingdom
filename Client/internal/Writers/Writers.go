package Writers

import (
	ConnClosed "kli/internal/ConnClosed"
	Read "kli/internal/reader"
	"log"
	"net"
)

type Request struct {
	To      string
	Command string
}

type Writer struct {
	Conn     net.Conn
	From     string
	Request  chan Request
	Response chan []byte
}

func WriterHandler(writer *Writer) {
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
				_, err = conn.Write([]byte(writer.From + " " + "!" + input.Command))
			} else {
				_, err = conn.Write([]byte(writer.From + " " + input.To + " " + input.Command))
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
