package main

import (
	zmq "github.com/pebbe/zmq4"
)

func main() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5563")
	subscriber.SetSubscribe("alerts")
	for {
		address, _ := subscriber.Recv(0)
		content, _ := subscriber.Recv(0)
		print("[" + string(address) + "] " + string(content) + "\n")
	}
}
