package listener

import (
	"fmt"
	notifiers "kingdom/internal/notifiers"
	tellers "kingdom/internal/tellers"
	"log"
	"net"
)

type Listener struct {
	Addr string
	Port int
}

func StartListener(address string, port int, callback chan<- tellers.Teller) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on %s:%d\n", address, port)
	num := 0

	name := "Fish_" + fmt.Sprint(port) + "_" + fmt.Sprint(num)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		newInput := make(chan []byte)
		newOutput := make(chan []byte)

		go tellers.StartTeller(conn, newInput, newOutput)
		ipAddress := conn.RemoteAddr().String()
		ipAddress = "New connection from: " + ipAddress

		fmt.Println(ipAddress)

		go notifiers.StartNotifier(callback, tellers.Teller{Name: name, Input: newInput, Output: newOutput})
		num += 1
		name = "Fish_" + fmt.Sprint(port) + "_" + fmt.Sprint(num)

	}
}

func Run(listener Listener, callback chan<- tellers.Teller) {
	go StartListener(listener.Addr, listener.Port, callback)
}
