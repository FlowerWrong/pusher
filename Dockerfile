FROM golang:latest

WORKDIR $GOPATH/bin
COPY ./pusher $GOPATH/bin
COPY ./settings.yml $GOPATH/bin

ENV APP_ENV=production

EXPOSE 8100
ENTRYPOINT ["./pusher"]
