package teller

import (
	"fmt"
	"log"
	"net"
)

type Teller struct {
	Input  chan []byte
	Output chan []byte
	Owner  string
}

func clearSlice(slice []byte) {
	for i := 0; i < len(slice); i += 1 {
		slice[i] = 0
	}
}

func TellerReader(conn net.Conn, output chan<- []byte) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}

		copied := make([]byte, n)
		copy(copied, buffer)

		output <- copied
		clearSlice(buffer)
	}
}

func StartTeller(conn net.Conn, input <-chan []byte, output chan<- []byte) {
	defer conn.Close()

	reader := make(chan []byte)
	go TellerReader(conn, reader)

	for {
		select {
		case cmd := <-input:
			_, err := conn.Write(cmd)
			if err != nil {
				fmt.Println(err)
			}
		case data := <-reader:
			output <- data
		}

		// send command ls to the client
		// command := []byte("ls\n")
		// _, err = conn.Write(command)
		// if err != nil {
		//	  log.Println(err)
		// 	  return
		// }
	}
}
