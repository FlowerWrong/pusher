FROM golang:latest

WORKDIR $GOPATH/src/github.com/FlowerWrong/pusher
COPY . $GOPATH/src/github.com/FlowerWrong/pusher
RUN go get -u -v ./...
RUN go build -o pusher cmd/main.go

EXPOSE 8100
ENTRYPOINT ["./pusher"]
