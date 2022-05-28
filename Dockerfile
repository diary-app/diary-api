#!/bin/sh
FROM golang:1.18 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . /app
RUN cd /app/cmd/diary_api_server && GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-w -extldflags "-static"' -o /app/diary-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/diary-api /app
COPY --from=build /app/internal/config/config.yml /app/internal/config/

ENTRYPOINT ["/app/diary-api"]