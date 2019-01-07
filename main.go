package main

import (
	"encoding/json"
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"net/http"
	"os"
)

//HookMessage is the message we receive from Alertmanager
type HookMessage struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

//Alert is a single alert.
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"EndsAt,omitempty"`
}

var addr = flag.String("addr", ":9098", "address to listen for webhook")
var pubaddr = flag.String("publisher", "tcp://*:5563", "address fot the publish socket")
var topic = flag.String("topic", "alerts", "default zmq topic to publish hook messages")
var endpoint = flag.String("endpoint", "/alerts", "default http endpoint for alertmanager")
var pubsub = make(chan HookMessage)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func SendHookMessageToPublisher(topic string) {
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()
	publisher.Bind(*pubaddr)

	for {
		select {
		case hookMessage := <-pubsub:
			{
				encoded, err := json.Marshal(hookMessage)
				if err != nil {
					// TODO: add a error log
					continue
				}
				publisher.Send(topic, zmq.SNDMORE)
				publisher.Send(string(encoded), 0)
			}
		}
	}
}

func ReceiveHookMessageFromAlertManager(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		{
			var message HookMessage
			err := json.NewDecoder(request.Body).Decode(&message)
			if err != nil {
				http.Error(writer, err.Error(), 400)
				return
			}
			defer request.Body.Close()
			pubsub <- message
		}
	default:
		http.Error(writer, "unsupported HTTP method", 400)
	}
}

func main() {
	flag.Parse()
	go SendHookMessageToPublisher(*topic)
	go http.HandleFunc(*endpoint, ReceiveHookMessageFromAlertManager)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
