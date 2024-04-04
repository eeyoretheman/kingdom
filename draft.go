package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"syscall"
)

type Request struct {
	From    string
	Target  string
	Command string
	Body    []byte
}

type Response struct {
	Target string
	Body   []byte
}

type Client struct {
	Input  chan []byte
	Output chan []byte
}

type Teller struct {
	Input  chan []byte
	Output chan []byte
	Owner  string
}

func ConnClosed(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	}

	return false
}

func Reader(conn net.Conn, callback chan<- []byte, errChan chan<- error) {
	buf := make([]byte, 16384)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			errChan <- err

			if ConnClosed(err) {
				return
			}
		}

		copied := make([]byte, n)
		copy(copied, buf)

		callback <- copied

		for i := 0; i < len(buf); i += 1 {
			buf[i] = 0
		}
	}
}

func TellerHandler(conn net.Conn, responses chan<- Response, teller *Teller) {
	reader := make(chan []byte)
	readerErr := make(chan error)

	go Reader(conn, reader, readerErr)

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

			if ConnClosed(err) {
				responses <- Response{"!", []byte(teller.Owner)}
			}
		}
	}
}

func ClientHandler(conn net.Conn, requests chan<- Request, client Client) {
	reader := make(chan []byte)
	readerErr := make(chan error)

	go Reader(conn, reader, readerErr)

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

			if ConnClosed(err) {
				return
			}
		}
	}
}

func Notifier[T *Teller | Client](object T, channel chan T) {
	channel <- object
}

func TellerListener(bind string, callback chan<- Response, tellerChan chan *Teller) {
	listener, err := net.Listen("tcp", bind)

	if err != nil {
		log.Printf("Could not bind on %s; Error: %s\n", bind, err)
		return
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Could not accept AGENT connection; Error: %s\n", err)
		}

		teller := Teller{Input: make(chan []byte), Output: make(chan []byte), Owner: "!"}

		go TellerHandler(conn, callback, &teller)
		go Notifier(&teller, tellerChan)
	}
}

func ClientListener(bind string, callback chan<- Request, clientChan chan Client) {
	listener, err := net.Listen("tcp", bind)

	if err != nil {
		log.Printf("Could not bind on %s; Error: %s\n", bind, err)
		return
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Could not accept CLIENT connection; Error: %s\n", err)
		}

		client := Client{Input: make(chan []byte), Output: make(chan []byte)}

		go ClientHandler(conn, callback, client)
		go Notifier(client, clientChan)
	}
}

func main() {
	clients := make(map[string]Client)
	tellers := make(map[string]*Teller)

	clientChannel := make(chan Client)
	tellerChannel := make(chan *Teller)

	requests := make(chan Request)
	responses := make(chan Response)

	go ClientListener("localhost:2222", requests, clientChannel)

	for {
		select {
		case request := <-requests:
			if request.Target == "!" {
				switch request.Command {
				case "lst":
					for name := range tellers {
						clients[request.From].Input <- []byte(fmt.Sprintf("%s, %s\n", name, tellers[name].Owner))
					}
				case "lsc":
					for name := range clients {
						clients[request.From].Input <- []byte(fmt.Sprintf("%s\n", name))
					}
				case "rmt":
					id := strings.TrimSuffix(string(request.Body), "\n")
					_, ok := tellers[id]

					if !ok {
						clients[request.From].Input <- []byte("No such teller.\n")
						break
					}

					delete(tellers, id)
				case "rmc":
					id := strings.TrimSuffix(string(request.Body), "\n")
					_, ok := clients[id]

					if !ok {
						clients[request.From].Input <- []byte("No such client.\n")
						break
					}

					delete(clients, id)
				case "tl":
					bind := strings.TrimSuffix(string(request.Body), "\n")
					go TellerListener(bind, responses, tellerChannel)
				case "cl":
					bind := strings.TrimSuffix(string(request.Body), "\n")
					go ClientListener(bind, requests, clientChannel)
				}
				break
			}

			switch request.Command {
			case "lock":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != "!" {
					clients[request.From].Input <- []byte("Already locked.\n")
					break
				}

				teller.Owner = request.From
				tellers[request.Target] = teller
			case "unlock":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != request.From {
					clients[request.From].Input <- []byte("You do not own that.\n")
					break
				}

				teller.Owner = "!"
				tellers[request.Target] = teller
			case "send":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != request.From {
					clients[request.From].Input <- []byte("You do not own that.\n")
					break
				}

				teller.Input <- request.Body
			}
		case response := <-responses:
			if response.Target == "!" {
				_, ok := clients[string(response.Body)]

				if !ok {
					log.Println("Orphaned response:", response)
					break
				}

				clients[string(response.Body)].Input <- []byte("Your teller died.\n")
				break
			}

			_, ok := clients[response.Target]

			if !ok {
				log.Println("Orphaned response:", response)
				break
			}

			clients[response.Target].Input <- response.Body
		case client := <-clientChannel:
			name := fmt.Sprint(rand.Int())
			clients[name] = client
			client.Input <- []byte(name + "\n")
		case teller := <-tellerChannel:
			name := fmt.Sprint(rand.Int())
			tellers[name] = teller
		}
	}
}
