# 1.17-alpine bug : standard_init_linux.go:228: exec user process caused: no such file or directory
ARG GOLANG_VERSION=1.17-alpine

# Building custom health checker
FROM golang:$GOLANG_VERSION as health-build-env

# Copying source
WORKDIR /go/src/app
COPY ./healthcheck /go/src/app

# Installing dependencies
RUN go get -d -v ./...

# Compiling
RUN go build -o /go/bin/healthchecker

# Building bouncer
FROM golang:$GOLANG_VERSION as build-env

# Copying source
WORKDIR /go/src/app
COPY ./bouncer /go/src/app

# Installing dependencies
RUN go get -d -v ./...

# Compiling
RUN go build -o /go/bin/app

FROM ghcr.io/linuxserver/baseimage-alpine:3.17

COPY --from=health-build-env /go/bin/healthchecker /app
COPY --from=build-env /go/bin/app /app

COPY /root /

EXPOSE 8080
