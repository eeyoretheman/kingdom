package main

import (
	"fmt"
	listeners "kingdom/internal/listeners"
	tellers "kingdom/internal/tellers"
)

// run listeners when a menu option is selected with the address and port
func printMenu() {
	fmt.Println("1. Start listener")
	fmt.Println("2. Exit")
	fmt.Println("3. Manual connect to listener")
	fmt.Print("Enter your choice: ")
}

// CLI + Backend all in 1 for now; separate later
func main() {
	printMenu()

	var tellerSlice []tellers.Teller
	var active tellers.Teller
	callback := make(chan tellers.Teller)

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		listener := listeners.Listener{Addr: "127.0.0.1", Port: 1337}
		listeners.Run(listener, callback)

		for {
			select {
			case t := <-callback:
				fmt.Printf("Added %s\n", t.Name)
				tellerSlice = append(tellerSlice, t)
				active = t // for testing
			case data := <-active.Output:
				fmt.Printf("Input from '%s': %s", active.Name, data)
				active.Input <- []byte("ls\n")
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
