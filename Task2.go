package main

import (
	"context"
	"log"
	"net"

	"github.com/yourusername/subpub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pubsub subpub.SubPub
}

func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	sub, err := s.pubsub.Subscribe(req.Key, func(msg interface{}) {
		if err := stream.Send(&pb.Event{Data: msg.(string)}); err != nil {
			log.Printf("Stream error: %v", err)
		}
	})
	if err != nil {
		return status.Error(codes.Internal, "Subscription failed")
	}

	<-stream.Context().Done()
	sub.Unsubscribe()
	return nil
}

func (s *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	if err := s.pubsub.Publish(req.Key, req.Data); err != nil {
		return nil, status.Error(codes.Internal, "Publish failed")
	}
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, &Server{pubsub: subpub.NewSubPub()})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	s.GracefulStop()
}
