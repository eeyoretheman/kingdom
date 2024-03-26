package notifier

import (
	tellers "kingdom/internal/tellers"
)

func StartNotifier(callback chan<- tellers.Teller, teller tellers.Teller) {
	callback <- teller
}
