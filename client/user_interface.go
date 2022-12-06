package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (app *Application) addMessage(msg Message) {
	app.Messages = append(app.Messages, msg)
}

func (app *Application) setupUI() (*widget.List, *widget.Entry) {
	messages := widget.NewList(
		func() int {
			return len(app.Messages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Messages")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%v: %v", app.Messages[i].client, app.Messages[i].body))
		})

	//messages.Resize(fyne.NewSize(800, 400))
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your message...")
	//entry.Resize(fyne.NewSize(3, 100))

	app.MessagesWidget = messages
	app.EntryWidget = entry

	return messages, entry
}
