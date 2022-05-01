package main

import (
	"log"
	"net"

	"grpc-server/pkg/pb"
	"grpc-server/pkg/service"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

const (
	NatsConnection = "nats://nats-server:4222"
)

var (
	svc *service.Server
	nc  *nats.Conn
)

func init() {
	var err error

	nc, err = nats.Connect(NatsConnection)
	if err != nil {
		log.Panic(err)
	}

	svc, err = service.New(nc)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	defer func() {
		nc.Drain()
		nc.Close()
	}()

	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panic(err)
	}

	log.Println("listening :50051")

	server := grpc.NewServer()
	pb.RegisterMessageServiceServer(server, svc)
	log.Panic(server.Serve(list))
}
