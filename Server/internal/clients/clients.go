package client

import (
	"fmt"
	notifiers "kingdom/internal/notifiers"
	"log"
	"net"
)

type Client struct {
	Input  chan notifiers.Message
	Output chan notifiers.Message
}

func clearSlice(slice []byte) {
	for i := 0; i < len(slice); i += 1 {
		slice[i] = 0
	}
}

func clientReader(conn net.Conn, output chan<- []byte) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("[CLIENTS] Could not read; Connection was probably closed.")
		}

		copied := make([]byte, n)
		copy(copied, buf)

		output <- buf

		clearSlice(buf)
	}
}

func EphemeralReader(conn net.Conn, callback chan<- []byte) {
	buf := make([]byte, 4096)

	n, err := conn.Read(buf)

	if err != nil {
		log.Print(err)
	}

	callback <- buf[:n]
}

func goClient(client Client, conn net.Conn) {
	reader := make(chan []byte)

	for {
		go EphemeralReader(conn, reader)
	}
}

func New(client Client, conn net.Conn) {
	go goClient(client, conn)
}
