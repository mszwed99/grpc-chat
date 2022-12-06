package main

import (
	"context"
	chatpb "gRPC-Chat/chat-protos"
	"google.golang.org/grpc"
	"io"
	"log"
)

func ConnectAndListen(ip string) {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(ip, opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("Cannot close connection")
			return
		}
	}(cc)

	c := chatpb.NewChatServiceClient(cc)
	stream, err := c.ConnectToStream(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to stream: %v", err)
	}

	connect := NewClientConnect(*name, stream)

	go func() {
		err := connect.receiveMessage()
		if err != nil {
			log.Fatalf("Could not reveive message to server: %v", err)
		}
	}()
	// go routines
	go func() {
		err := connect.sendMessage()
		if err != nil {
			log.Fatalf("Could not send message to server: %v", err)
		}
	}()

	// Wait chan
	wait := make(chan bool)
	<-wait
}
func NewClientConnect(name string, stream chatpb.ChatService_ConnectToStreamClient) ClientConnect {
	return ClientConnect{
		name:   name,
		stream: stream,
	}
}

func (cc *ClientConnect) sendMessage() error {
	for msg := range msgchan {
		ccMessage := &chatpb.MessageRequest{
			Message: &chatpb.Message{
				Client: &chatpb.Client{
					Name: cc.name,
				},
				Message: msg,
			},
		}
		err := cc.stream.Send(ccMessage)
		if err != nil {
			log.Fatalf("Error while sending message to server: %v", err)
			return err
		}
	}
	return nil
}

func (cc *ClientConnect) receiveMessage() error {
	for {
		msg, err := cc.stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving message form server: %v", err)
			return err
		}

		m := Message{
			body:   msg.Message.GetMessage(),
			client: msg.Message.GetClient().GetName(),
		}

		application.addMessage(m)
		application.MessagesWidget.Refresh()
		application.MessagesWidget.ScrollToBottom()
		application.MessagesWidget.Refresh()
	}
}
