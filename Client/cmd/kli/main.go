package main

import (
	// Consider termui instead
	reader "kli/internal/reader"
	"log"
	"net"
)

type Response struct {
	From    string
	To      string
	Command string
	Body    []byte
}

type Request struct {
	From    string
	To      string
	Command string
}

func connect_to_main_server(request chan Request, response chan Response) {
	conn, err := net.Dial("tcp", "localhost:2222")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	read := make(chan []byte)
	errChan := make(chan error)
	go reader.Reader(conn, read, errChan)

	r := make(chan Request)

	for {
		select {
		case req := <-request:
			if req.To == "!" {
				conn.Write([]byte("! ! " + req.Command))
			} else {
				conn.Write([]byte(req.From + " " + req.To + " " + req.Command))
			}
			r <- req
		case data := <-read:
			req := <-r
			response <- Response{From: req.From, To: req.To, Command: req.Command, Body: data}
		case err := <-errChan:
			log.Fatal(err)
		}
	}
}

func main() {
	request := make(chan Request)
	response := make(chan Response)

	go connect_to_main_server(request, response)
	// ill fix it later probably this day but will push for now
}
