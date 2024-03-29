FROM golang:1.18 AS builder

WORKDIR /app/delivery

COPY . /app/delivery

RUN mkdir bin

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/cli ./cmd/cli/main.go


FROM alpine:latest as prod

WORKDIR /app/delivery

RUN mkdir bin
RUN mkdir migrations
RUN mkdir check
RUN mkdir logs
RUN touch logs/log.txt

COPY --from=builder /app/delivery/bin ./bin
COPY --from=builder /app/delivery/migrations ./migrations
COPY --from=builder /app/delivery/check ./check
COPY --from=builder /app/delivery/prod.config.yaml .

CMD ["sh","-c", "bin/app"]