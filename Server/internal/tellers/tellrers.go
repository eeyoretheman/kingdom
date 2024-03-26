package teller

import (
	"fmt"
	"log"
	"net"
)

func Teller(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err :=
			conn.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		data := buffer[:n]
		fmt.Printf("Received data: %s\n", string(data))

		// send command ls to the client
		command := []byte("ls\n")
		_, err = conn.Write(command)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
