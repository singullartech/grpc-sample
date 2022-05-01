package grpc

import (
	"context"
	"fmt"
	"grpc-stream/pkg/pb"
	"log"

	"google.golang.org/grpc"
)

type MessageStream struct {
	client   pb.MessageServiceClient
	messages int
}

func NewMessageStream(conn grpc.ClientConnInterface, q int) MessageStream {
	return MessageStream{
		client:   pb.NewMessageServiceClient(conn),
		messages: q,
	}
}

func (s MessageStream) Send(ctx context.Context, pool int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.client.Send(ctx)
	if err != nil {
		return fmt.Errorf("create stream: %w", err)
	}

	total := s.messages

	for i := 1; i <= total; i++ {
		req := &pb.MessageRequest{
			Message: &pb.Message{
				Body: fmt.Sprintf("Hello thread %d!!!", pool),
			},
		}

		if err := stream.Send(req); err != nil {
			log.Fatalln("send stream: %w", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close and receive: %w", err)
	}

	log.Printf("routine %03d %+v\n", pool, response)
	return nil
}
