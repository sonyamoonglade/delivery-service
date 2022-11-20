FROM golang:1.18

ARG APP_NAME=delivery-service

WORKDIR /app

COPY . .

RUN apt update -y && \
    apt upgrade -y && \
    apt install -y git && \
    go install github.com/githubnemo/CompileDaemon@latest


ENTRYPOINT CompileDaemon -polling -build="go build -o ./bin/${APP_NAME} ./cmd/app/main.go && go build -o ./bin/cli ./cmd/cli/main.go" -command="./bin/${APP_NAME} --debug --strict=false"
