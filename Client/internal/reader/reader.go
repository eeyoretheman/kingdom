package reader

import (
	ConnClosed "kli/internal/ConnClosed"
	"net"
)

func Reader(conn net.Conn, callback chan<- []byte, errChan chan<- error) {
	buf := make([]byte, 16384)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			errChan <- err

			if ConnClosed.ConnClosed(err) {
				return
			}
		}

		copied := make([]byte, n)
		copy(copied, buf)

		callback <- copied

		for i := 0; i < len(buf); i += 1 {
			buf[i] = 0
		}
	}
}
