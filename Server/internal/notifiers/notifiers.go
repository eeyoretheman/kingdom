package notifiers

import (
	clients "kingdom/internal/clients"
	tellers "kingdom/internal/tellers"
)

func Notifier[T *tellers.Teller | clients.Client](object T, channel chan T) {
	channel <- object
}
