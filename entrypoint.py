#!/usr/bin/env python3
import concurrent.futures
import time
import socket
import sys

socket.setdefaulttimeout(1)

def lookup_host():
    print(socket.gethostbyname(sys.argv[1]))

with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
    while True:
        executor.submit(lookup_host)
        time.sleep(0.05)
