# Pusher

* [pusher-websockets-protocol reference](https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol)
* [rest-api reference](https://pusher.com/docs/channels/library_auth_reference/rest-api)

## How to start?

```bash
go run cmd/main.go --config=./settings.yml

GOOS=linux GOARCH=amd64 go build -o pusher cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o pusher cmd/main.go
```

## Requirements

* [golang](https://golang.org/)
* [redis](https://redis.io/): storage engine and pubsub

## Features

* [x] public channels
* [x] private channels
* [x] presence channels
* [x] client events
* [x] rest api
* [x] webhooks
* [ ] encrypted channels(BETA)
* [x] distributed
* [ ] integration test
* [ ] benchmark
* [ ] redis master slave
* [ ] support other pubsub message brokers
* [ ] apm

## Other implementations

* [slanger](https://github.com/stevegraham/slanger)
