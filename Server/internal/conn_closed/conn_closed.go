package ConnClosed

import (
	"errors"
	"io"
	"net"
	"syscall"
)

func ConnClosed(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	}

	return false
}
