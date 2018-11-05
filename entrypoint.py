#!/usr/bin/env python3
import concurrent.futures
import time
import socket
import sys

socket.setdefaulttimeout(1)

hostname = sys.argv[1]

def lookup_host(name):
    print(socket.getaddrinfo(name, None))

with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
    while True:
        executor.submit(lookup_host, hostname)
        time.sleep(0.05)
