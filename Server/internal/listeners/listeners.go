package listeners

import (
	"fmt"
	"log"
	"net"
)

func startListener(address string, port int) {
	// Create a TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on %s:%d\n", address, port)

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from the connection
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	// Process the received data
	data := buffer[:n]
	fmt.Printf("Received data: %s\n", string(data))

	// Send a response back to the client
	response := []byte("Hello from the server!")
	_, err = conn.Write(response)
	if err != nil {
		log.Println(err)
		return
	}
}

func run(listeners []struct {
	address string
	port    int
}) {
	// listeners := []struct {
	// 	address string
	// 	port    int
	// }{
	// 	{"0.0.0.0", 1337},
	// 	{"0.0.0.0", 1338},
	// }

	for _, l := range listeners {
		go startListener(l.address, l.port)
	}

	select {}
}
