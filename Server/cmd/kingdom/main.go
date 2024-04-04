package main

import (
	"fmt"
	clients "kingdom/internal/clients"
	listeners "kingdom/internal/listeners"
	notifiers "kingdom/internal/notifiers"
	tellers "kingdom/internal/tellers"
	"net"
	"strconv"
	"strings"
)

// Default client port: 2222
// NOTE: Closing the client (netcat) crashes the server; Make the server wait for a new connection if the old one closed.

// run listeners when a menu option is selected with the address and port
func printMenu() {
	fmt.Println("1. Start listener")
	fmt.Println("2. Exit")
	fmt.Println("3. Manual connect to listener")
	fmt.Print("Enter your choice: ")
}

func debugReader(conn net.Conn, output chan<- []byte) {
	for {
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		if err != nil {
			panic("[DEBUG] Could not read.")
		}

		copied := make([]byte, n)
		copy(copied, buf)

		output <- copied
	}
}

// Refactor later
func debugHandler(channel chan []byte) {
	listener, err := net.Listen("tcp", "localhost:2222")

	if err != nil {
		panic("[DEBUG] Could not bind on :2222.")
	}

	conn, err := listener.Accept()

	if err != nil {
		panic("[DEBUG] Could not accept connection.")
	}

	input := make(chan []byte)

	for {
		go debugReader(conn, input)

		select {
		case data := <-channel:
			conn.Write(data)
		case bytes := <-input:
			channel <- bytes
		}
	}
}

// CLI + Backend all in 1 for now; separate later
func old_main() {
	printMenu()

	var tellerList []tellers.Teller
	var active tellers.Teller
	callback := make(chan tellers.Teller)

	clientDebug := make(chan []byte)

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		listener := listeners.Listener{Addr: "127.0.0.1", Port: 1337}
		listeners.New(listener, callback)

		go debugHandler(clientDebug)

		for {
			select {
			case t := <-callback:
				fmt.Printf("Added a new agent.\n")
				tellerList = append(tellerList, t)
				active = t // Set the newly added agent as active; Consider removing
			case data := <-active.Output:
				fmt.Printf("Received: %s", data)
				active.Input <- []byte("ls\n")
			case input := <-clientDebug:
				p := strings.Split(strings.TrimSuffix(string(input), "\n"), " ")

				fmt.Println(p)

				switch p[0] {
				case "active":
					fmt.Println(len(p))
					if len(p) < 2 {
						clientDebug <- []byte("Not enough args.\n")
						break
					}
					num, _ := strconv.Atoi(p[1])
					active = tellerList[num]
				case "ls":
					for i := range tellerList {
						clientDebug <- []byte(fmt.Sprintf("%d %s\n", i))
					}
				// NOTE: Allow for quoted strings and escapes
				case "send":
					active.Input <- []byte(p[1])
				default:
					clientDebug <- []byte("No such command.\n")
				}
			}
		}
	case 2:
		fmt.Println("Exiting...")

	case 3:
		// connect to the listener with local connection
		// conn, err := net.Dial("tcp", "127.0.0.1:1337")
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
}

type Message struct {
	To   string
	From string
	Body []byte
}

func main() {
	clientList := make(map[string]clients.Client)
	tellerList := make(map[string]tellers.Teller)

	interpreterInput := make(chan []byte)
	interpreterOutput := make(chan []byte)

	clientChannel := make(chan notifiers.Message)
	tellerChannel := make(chan notifiers.Message)

	clients.New(clients.Client{})

	for {
		select {
		case message := <-clientChannel:
			if message.To == "!" {
				cmd := strings.Split(strings.TrimSuffix(string(message.Body), "\n"), " ")
				switch cmd[0] {
				case "ls":
					for i, e := range tellerList {
						clientList[message.From].Input <- notifiers.Message{From: "!", Body: []byte(fmt.Sprintf("%d %s", i, e))}
					}
				case "active":
					continue
				default:
				}
			} else {

			}
		case message := <-tellerChannel:
			fmt.Println(message, tellers)
		}
	}
}
