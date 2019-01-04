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

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {
	pubsub := make(chan HookMessage)

	addr := flag.String("addr", ":9098", "address to listen for webhook")
	pubaddr := flag.String("publisher", "tcp://*:5563", "address fot the publish socket")

	flag.Parse()

	go func(ch chan<- HookMessage) {

		go func(ch <-chan HookMessage, topic string) {

			publisher, _ := zmq.NewSocket(zmq.PUB)
			defer publisher.Close()
			publisher.Bind(*pubaddr)

			for {
				select {
				case message := <-pubsub:
					{
						encoded, err := json.Marshal(message)
						if err != nil {
							fmt.Println("failed")
							continue
						}
						publisher.Send(topic, zmq.SNDMORE)
						publisher.Send(string(encoded), 0)
					}
				}
			}

		}(pubsub, "alerts")

		http.HandleFunc("/alerts", func(writer http.ResponseWriter, request *http.Request) {
			switch request.Method {
			case http.MethodPost:
				{
					dec := json.NewDecoder(request.Body)
					defer request.Body.Close()
					var message HookMessage
					if err := dec.Decode(&message); err != nil {
						log.Printf("error decoding message: %v", err)
						http.Error(writer, "invalid request body", 400)
						return
					}
					ch <- message
				}
			default:
				http.Error(writer, "unsupported HTTP method", 400)
			}
		})

	}(pubsub)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
