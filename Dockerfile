# Building bouncer
FROM docker.io/golang:1.17-alpine as build-env

# Copying source
WORKDIR /go/src/app
COPY ./bouncer /go/src/app

# Installing dependencies
RUN go get -d -v ./...

# Compiling
RUN go build -o /go/bin/bouncer

FROM ghcr.io/linuxserver/baseimage-alpine:3.17

COPY --from=build-env /go/bin/bouncer /app

COPY /root /

EXPOSE 8080
