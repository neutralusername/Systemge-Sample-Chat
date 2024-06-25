package appChat

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.JOIN:  app.Join,
		topics.LEAVE: app.Leave,
	}
}

func (app *App) Join(client *Client.Client, message *Message.Message) (string, error) {
	if err := app.AddChatter(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to create chatter", err)
	}
	if err := app.AddToRoom(message.GetOrigin(), message.GetPayload()); err != nil {
		return "", Utilities.NewError("Failed to join room", err)
	}
	return Utilities.StringsToJsonObjectArray(app.GetRoomMessages(message.GetPayload())), nil
}

func (app *App) Leave(client *Client.Client, message *Message.Message) (string, error) {
	if err := app.RemoveFromRoom(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to leave room", err)
	}
	if err := app.RemoveChatter(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to leave room", err)
	}
	return "", nil
}
