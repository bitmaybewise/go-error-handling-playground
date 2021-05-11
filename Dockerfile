FROM golang:1.16 AS base
WORKDIR /app
COPY . .
# adding reflex to add hot reloading capability
RUN go get github.com/cespare/reflex
RUN go mod tidy
