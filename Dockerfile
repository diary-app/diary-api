FROM golang:1.17-alpine as build-env
WORKDIR /app
ADD . /app
RUN cd /app/cmd/diary_api_server && go build -o /app/diary-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/diary-api /app

ENTRYPOINT ["/diary-api"]