package main

import (
	"fmt"
	chatpb "gRPC-Chat/chat-protos"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("gRPC Chat Server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, &server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to run a gRPC server: %v", err)
	}

}
