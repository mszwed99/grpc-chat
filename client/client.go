package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var application Application
var msgchan chan string

var (
	address *string
	port    *string
	name    *string
)

func handleMessage(msg string) {
	msgchan <- msg
	application.EntryWidget.Text = ""
	application.EntryWidget.Refresh()
	application.MessagesWidget.ScrollToBottom()
	application.MessagesWidget.Refresh()
}

func main() {
	name = flag.String("n", "name", "Client name")
	address = flag.String("s", "localhost", "gRPC server address")
	port = flag.String("p", ":50051", "gRPC server port")
	flag.Parse()
	ip := *address + *port

	msgchan = make(chan string)

	go ConnectAndListen(ip)

	//UI
	a := app.New()
	win := a.NewWindow(fmt.Sprintf("gRPC Chat Client - [%v]", *name))

	messages, entry := application.setupUI()
	entry.OnSubmitted = func(s string) {
		handleMessage(entry.Text)
	}
	//btn := widget.NewButton("Send", func() {
	//	handleMessage(entry.Text)
	//})

	//entryBtnContainer := container.NewHSplit(entry, btn)
	content := container.NewVSplit(messages, entry)
	content.Offset = 2
	win.SetContent(content)

	// Show windows and run app
	win.Resize(fyne.Size{Width: 900, Height: 600})
	win.CenterOnScreen()
	win.ShowAndRun()

}
