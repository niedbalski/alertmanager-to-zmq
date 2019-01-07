# Prometheus Alertmanager to ZMQ [![Build Status][buildstatus]][circleci] [![Docker Repository on Quay](https://quay.io/repository/niedbalski/openstack-exporter/status "Docker Repository on Quay")](https://quay.io/repository/niedbalski/openstack-exporter)

This program receives HTTP messages from a prometheus alertmanager and 
proxies those messages to a ZMQ socket.

```sh
./alertmanager-to-zmq --help
Usage of ./alertmanager-to-zmq:
  -addr string
        address to listen for webhook (default ":9098")
  -endpoint string
        default http endpoint for alertmanager (default "/alerts")
  -publisher string
        address fot the publish socket (default "tcp://*:5563")
  -topic string
        default zmq topic to publish hook messages (default "alerts")
```

Example prometheus alertmanager configuration:


```yaml
route:
  receiver: webhook
  group_wait: 0s
  group_interval: 1s
  repeat_interval: 1s

receivers:
- name: "webhook"
  webhook_configs:
  - url: http://localhost:9098/alerts
```

Check the example_client directory for further examples.


Alternatively a Dockerfile and image are supplied

```sh
docker run -p 9180:9180 quay.io/niedbalski/alertmanager-to-zmq:v0.0.1
```

[buildstatus]: https://circleci.com/gh/niedbalski/alertmanager-to-zmq/tree/master.svg?style=shield
[circleci]: https://circleci.com/gh/niedbalski/alertmanager-to-zmq
[hub]: https://hub.docker.com/r/niedbalski/alertmanager-to-zmq/
