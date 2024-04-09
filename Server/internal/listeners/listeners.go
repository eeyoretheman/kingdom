package listeners

import (
	clients "kingdom/internal/clients"
	notifiers "kingdom/internal/notifiers"
	tellers "kingdom/internal/tellers"
	"log"
	"net"
)

func TellerListener(bind string, callback chan<- tellers.Response, tellerChan chan *tellers.Teller) {
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

		teller := tellers.Teller{Input: make(chan []byte), Output: make(chan []byte), Owner: "!"}

		go tellers.TellerHandler(conn, callback, &teller)
		go notifiers.Notifier(&teller, tellerChan)
	}
}

func ClientListener(bind string, callback chan<- clients.Request, clientChan chan clients.Client) {
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

		client := clients.Client{Input: make(chan []byte), Output: make(chan []byte)}

		go clients.ClientHandler(conn, callback, client)
		go notifiers.Notifier(client, clientChan)
	}
}
