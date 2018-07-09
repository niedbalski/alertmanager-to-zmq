#!/bin/en python3
# pip install pyzmq

import zmq
import json

context = zmq.Context()
socket = context.socket(zmq.SUB)
socket.connect('tcp://localhost:5563')

socket.setsockopt(zmq.SUBSCRIBE, 'alerts')

while True:
    header = socket.recv()
    alert = json.loads(socket.recv())
    print("[%s] %s" % (header, alert))
