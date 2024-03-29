FROM golang:1.18 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . /app
RUN cd /app/cmd/diary_api_server && \
    GOOS=linux GOARCH=amd64 go build -o /app/diary-api

FROM alpine
RUN apk update && apk add ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=build /app/diary-api /app
COPY --from=build /app/internal/config/config.yml /app/internal/config/

ENTRYPOINT ["/app/diary-api"]