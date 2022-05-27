FROM golang:1.18-alpine as build-env
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . /app
RUN cd /app/cmd/diary_api_server && go build -o /app/diary-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/diary-api /app
RUN ls -la /app/*
COPY --from=build-env /app/internal/config/config.yml /app/internal/config/

ENTRYPOINT ["/app/diary-api"]