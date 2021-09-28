FROM golang:1.17.1-alpine AS build

WORKDIR /tmp/app

COPY . .

RUN apk add --no-cache git && \
    go mod download && \
    go mod verify && \
    go build -o app

FROM alpine:latest

WORKDIR /home/app

COPY --from=build /tmp/app/app /home/app/

EXPOSE 80

ENTRYPOINT ./app