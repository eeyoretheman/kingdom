package main

import (
	"fmt"
	listeners "kingdom/internal/listeners"
	tellers "kingdom/internal/tellers"
	"net"
	"strconv"
	"strings"
	"time"
)

var TellerSlice [][]tellers.Teller

func get_response(t tellers.Teller) {
	for {
		data := <-t.Output
		fmt.Printf("Data received from %s: %s", t.Name, string(data))
	}
}

func handleListener(listener listeners.Listener, callback chan tellers.Teller) {
	listeners.Run(listener, callback)
	for {
		select {
		case t := <-callback:
			fmt.Printf("Added %s\n", t.Name)
			TellerSlice[listener.Port] = append(TellerSlice[listener.Port], t)
			go get_response(t)
		//if 10 seconds have passed, send a command to the listener
		case <-time.After(10 * time.Second):
			for _, t := range TellerSlice[listener.Port] {
				command := []byte("ls\n")
				t.Input <- command
			}
		}
	}
}

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

func debugger(clientDebug chan []byte, active tellers.Teller) {
	for {
		input := <-clientDebug
		p := strings.Split(strings.TrimSuffix(string(input), "\n"), " ")

		fmt.Println(p)

		switch p[0] {
		case "active":
			fmt.Println(len(p))
			if len(p) < 3 {
				clientDebug <- []byte("Not enough args.\n")
				break
			}
			port, _ := strconv.Atoi(p[1])
			num, _ := strconv.Atoi(p[2])
			active = TellerSlice[port][num]
		case "ls":
			port, _ := strconv.Atoi(p[1])
			for i, e := range TellerSlice[port] {
				clientDebug <- []byte(fmt.Sprintf("%d %s\n", i, e.Name))
			}
		// NOTE: Allow for quoted strings and escapes
		case "send":
			active.Input <- []byte(p[1])
		default:
			clientDebug <- []byte("No such command.\n")
		}
	}
}

// CLI + Backend all in 1 for now; separate later
func main() {
	printMenu()

	var active tellers.Teller
	callback := make(chan tellers.Teller)
	listener := listeners.Listener{Addr: "127.0.0.1", Port: 1337}
	var choice int = 0
	clientDebug := make(chan []byte)

	for choice != 2 {
		printMenu()
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			go handleListener(listener, callback)
			go debugHandler(clientDebug)
			go debugger(clientDebug, active)
			listener.Port += 1

		case 2:
			fmt.Println("Exiting...")
			return

		case 3:
			// connect to the listener with local connection
			// conn, err := net.Dial("tcp", "127.0.0.1:1337")
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
}
