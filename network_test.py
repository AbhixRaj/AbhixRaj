#!/usr/bin/env python3

import sys
import os
from pythonping import ping
import socket

def ping_host(host):
    print(f"Pinging {host}...")
    response = ping(host, count=4, verbose=True)
    print(response)

def port_scan(host, port):
    print(f"Scanning {host}:{port}...")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(1)
    result = sock.connect_ex((host, port))
    if result == 0:
        print(f"Port {port} is OPEN on {host}.")
    else:
        print(f"Port {port} is CLOSED on {host}.")
    sock.close()

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python network_test.py <ping|scan> <host> [port]")
        sys.exit(1)

    command, host = sys.argv[1:3]
    if command == "ping":
        ping_host(host)
    elif command == "scan" and len(sys.argv) == 4:
        port = int(sys.argv[3])
        port_scan(host, port)
    else:
        print("Invalid command or arguments.")
