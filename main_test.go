package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

const testMessage = `
{
	"version": "4",
	"groupKey": "{}:{alertname=\"TestAlert\"}",
	"status": "firing",
	"receiver": "team",
	"groupLabels": {},
	"commonLabels": {},
	"commonAnnotations": {},
	"externalURL": "http://127.0.0.1:9093",
	"alerts": [{
		"status": "firing",
		"labels": {},
		"annotations": {},
		"startsAt": "2019-01-03T22:20:08.822+09:00",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://127.0.0.1:9090/graph"
	}]
}
`

func TestSendHookMessageToPublisher(t *testing.T) {
	wg := sync.WaitGroup{}
	go SendHookMessageToPublisher(*topic)

	wg.Add(1)
	go func() {
		var received HookMessage
		var sent HookMessage

		defer wg.Done()
		subscriber, _ := zmq.NewSocket(zmq.SUB)
		subscriber.Connect("tcp://localhost:5563")
		subscriber.SetSubscribe(*topic)
		defer subscriber.Close()

		address, err := subscriber.Recv(0)
		if err != nil {
			panic(err)
		}

		content, err := subscriber.Recv(0)
		if err != nil {
			panic(err)
		}

		assert.NotNil(t, address)
		assert.NotNil(t, content)

		_ = json.Unmarshal([]byte(content), &sent)
		_ = json.Unmarshal([]byte(testMessage), &received)

		assert.Equal(t, received, sent)
	}()

	router := mux.NewRouter()
	router.HandleFunc("/alerts", ReceiveHookMessageFromAlertManager)
	req, _ := http.NewRequest("POST", "/alerts",
		bytes.NewBuffer([]byte(testMessage)))

	req.Header.Set("Content-Type", "application/json")

	<-time.After(1 * time.Second)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200)
	wg.Wait()
}
