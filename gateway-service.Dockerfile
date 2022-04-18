FROM golang:1.17 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o gateway ./src/gateway-service/cmd/gateway-service

ENTRYPOINT ["./gateway"]