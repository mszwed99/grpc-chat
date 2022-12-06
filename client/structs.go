package main

import (
	"fyne.io/fyne/v2/widget"
	chatpb "gRPC-Chat/chat-protos"
)

type Application struct {
	MessagesWidget *widget.List
	EntryWidget    *widget.Entry

	Messages []Message
}

type Message struct {
	body   string
	client string
}
type ClientConnect struct {
	stream chatpb.ChatService_ConnectToStreamClient
	name   string
}
