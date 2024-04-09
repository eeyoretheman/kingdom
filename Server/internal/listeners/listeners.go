package listener

import (
	clients "kingdom/internal/clients"
	"log"
	"net"
)

type Listener struct {
	Address string
	Port    int
}

func goClientListener(bind string, callback chan clients.Client) {
	listener, err := net.Listen("tcp", bind)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
		}

		go clients.New()
	}
}

func New(listener Listener) {

}
