FROM golang:1.14 as base

WORKDIR /app

RUN go get -u github.com/githubnemo/CompileDaemon
ENV GO111MODULE "on"
ENTRYPOINT CompileDaemon -build="go build -mod=vendor -o bin/chat_server ./cmd/..." -command="bin/chat_server start" -exclude-dir=vendor -exclude-dir=.git -exclude-dir=bin -graceful-kill
