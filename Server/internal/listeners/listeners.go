package listeners

import (
	"fmt"
	notifier "kingdom/internal/notifier"
	teller "kingdom/internal/tellers"
	"log"
	"net"
)

type Listener struct {
	Addr string
	Port int
}

func startListener(address string, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on %s:%d\n", address, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go teller.Teller(conn)
		ip_address := conn.RemoteAddr().String()
		ip_address = "New connection from: " + ip_address
		notifier.Notifier(ip_address)
	}
}

func Run(listen Listener) {

	go startListener(listen.Addr, listen.Port)

	select {}
}
