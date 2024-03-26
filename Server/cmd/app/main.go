package main

import (
	"fmt"
	"kingdom/internal/listeners"
)

// run listeners when a menu option is selected with the address and port
func print_menu() {
	fmt.Println("1. Start listener")
	fmt.Println("2. Exit")
	fmt.Println("3. Manual connect to listener")
	fmt.Print("Enter your choice: ")
}

func main() {
	print_menu()
	var choice int
	fmt.Scanln(&choice)
	switch choice {

	case 1:
		listen := listeners.Listener{Addr: "127.0.0.1", Port: 1337}
		listeners.Run(listen)
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
