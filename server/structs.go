package main

import (
	chatpb "gRPC-Chat/chat-protos"
	"sync"
)

type Server struct {
	chatpb.UnimplementedChatServiceServer
	messages []Message
	clients  map[string]chatpb.ChatService_ConnectToStreamServer
	mu       sync.RWMutex
}

type Client struct {
	id   string
	name string
}

type Message struct {
	id      string
	client  Client
	message string
}

type MessageQueue struct {
	mu  sync.Mutex
	que []Message
}
