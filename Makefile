client:
	go build example_client/main.go && mv main client
build:
	docker build -t alertmanager-to-zmq .
