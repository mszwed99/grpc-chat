package main

import (
	chatpb "gRPC-Chat/chat-protos"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

var server Server
var messageQueue MessageQueue

// Client functions

func (s *Server) addClient(uid string, stream chatpb.ChatService_ConnectToStreamServer) {
	if s.clients == nil {
		s.clients = make(map[string]chatpb.ChatService_ConnectToStreamServer)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[uid] = stream
}

func (s *Server) removeClient(uid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, uid)
}

func (s *Server) getClients() []chatpb.ChatService_ConnectToStreamServer {
	var clients []chatpb.ChatService_ConnectToStreamServer
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.clients {
		clients = append(clients, c)
	}

	return clients
}

func (s *Server) ConnectToStream(stream chatpb.ChatService_ConnectToStreamServer) error {
	uid := uuid.Must(uuid.NewRandom()).String()
	s.addClient(uid, stream)
	defer s.removeClient(uid)

	// Error channel
	errchan := make(chan error)

	// Broadcast every received message
	go s.ReceiveAndBroadcast(stream, uid, errchan)

	return <-errchan
}

func (s *Server) broadcast(msg Message) {
	// validate msg!!!!
	if len(msg.message) > 0 {
		m := &chatpb.MessageResponse{
			Message: &chatpb.Message{
				Client: &chatpb.Client{
					Name: msg.client.name,
				},
				Message: msg.message,
			},
		}

		for _, c := range s.getClients() {
			err := c.Send(m)
			if err != nil {
				log.Printf("Broadcast err: %v", err)
			}
		}
	}

}

func (s *Server) ReceiveAndBroadcast(stream chatpb.ChatService_ConnectToStreamServer, uid string, errchan chan error) {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			errchan <- nil
		}

		if err != nil {
			// Client disconnect error
			streamErr, ok := status.FromError(err)
			if ok {
				if streamErr.Code() == codes.Canceled {
					errchan <- nil
				}
			} else {
				log.Printf("Error while receiving message from client: %v", err)
				errchan <- err
			}

		}
		if msg != nil {
			m := Message{
				client: Client{
					id:   uid,
					name: msg.GetMessage().GetClient().GetName(),
				},
				message: msg.GetMessage().GetMessage(),
			}

			// Send received message to all clients
			messageQueue.mu.Lock()
			s.broadcast(m)
			messageQueue.mu.Unlock()
		}

	}

}
