FROM golang:latest

WORKDIR $GOPATH/bin
COPY ./pusher $GOPATH/bin
COPY ./settings.yml $GOPATH/bin

ENV APP_ENV=production BIND_ADDR=127.0.0.1 BIND_PORT=8100

EXPOSE 8100
ENTRYPOINT ["./pusher"]
