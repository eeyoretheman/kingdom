package notifier

import (
	tellers "kingdom/internal/tellers"
)

type Message struct {
	From string
	Body []byte
}

func StartNotifier(callback chan<- tellers.Teller, teller tellers.Teller) {
	callback <- teller
}
