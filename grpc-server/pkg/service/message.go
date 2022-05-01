package service

import (
	"encoding/json"
	"io"
	"log"

	"grpc-server/pkg/pb"

	"github.com/nats-io/nats.go"
)

type Server struct {
	nc *nats.Conn
}

const (
	MessageSubject = "sample.sub.messages"
)

func New(nc *nats.Conn) (*Server, error) {
	return &Server{nc: nc}, nil
}

func (s Server) Send(stream pb.MessageService_SendServer) error {
	var total int32

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&pb.MessageResponse{
					Success: true,
					Total:   total,
				})
		}

		data, err := json.Marshal(message.GetMessage())
		if err != nil {
			return err
		}

		err = s.nc.Publish(MessageSubject, data)
		if err != nil {
			log.Fatalln("error natsss ", err)
			return err
		}

		total++
	}
}
