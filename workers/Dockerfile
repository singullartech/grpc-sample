FROM golang:1.17-alpine AS base
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

FROM base AS builder
WORKDIR /app

COPY . /app
RUN go mod download \
    && go mod verify

RUN go build -o worker -a .

### Production
FROM alpine:3.14
COPY --from=builder /app/worker /usr/local/bin/worker

ARG USER=app
ENV HOME /home/$USER

WORKDIR /home/app

RUN adduser -D $USER \
    && chown $USER:$USER -R $HOME \
    && touch .env

USER $USER
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/worker"]