FROM golang:latest

WORKDIR $GOPATH/bin
COPY ./pusher $GOPATH/bin
COPY ./settings.yml $GOPATH/bin

ENV APP_ENV=production PORT=8100 REDIS_URL=redis://:@localhost:6379/1

EXPOSE 8100
ENTRYPOINT ["./pusher"]
