### base
FROM golang:1.18-alpine AS base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    build-base \
    git \
    && update-ca-certificates

COPY go.* ./

### builder
FROM base AS builder
WORKDIR /app

COPY . /app
RUN go mod download \
    && go mod verify

RUN go build -o grpc-server -a .

### production
FROM alpine:3.14
COPY --from=builder /app/grpc-server /usr/local/bin/grpc-server

ARG USER=app
ENV HOME /home/$USER

WORKDIR /home/app

RUN adduser -D $USER \
    && chown $USER:$USER -R $HOME \
    && touch .env

EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/grpc-server"]