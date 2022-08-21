FROM golang:1.18

WORKDIR /app/delivery

COPY . /app/delivery

RUN mkdir bin && \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/app/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/cli ./cmd/cli/main.go

CMD ["sh","-c","bin/app"]