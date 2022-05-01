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

FROM base AS dev
WORKDIR /app

RUN go get -u github.com/mgechev/revive@master \
    && go get -u google.golang.org/grpc \
    && go get -u google.golang.org/protobuf \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1 \
    && go mod tidy
    
RUN apk add --no-cache make protobuf-dev=3.18.1-r1    

RUN apk add curl \
    && curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin


### builder
FROM base AS builder
WORKDIR /app

COPY . /app
RUN go mod download \
    && go mod verify

RUN go build -o grpc-client -a .

### production
FROM alpine:3.14
COPY --from=builder /app/grpc-client /usr/local/bin/grpc-client

ARG USER=app
ENV HOME /home/$USER

WORKDIR /home/app

RUN adduser -D $USER \
    && chown $USER:$USER -R $HOME
