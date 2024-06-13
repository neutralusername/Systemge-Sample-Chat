package appChat

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.JOIN:  app.Join,
		topics.LEAVE: app.Leave,
	}
}

func (app *App) Join(message *Message.Message) (string, error) {
	chatter, err := app.AddChatter(message.GetOrigin())
	if err != nil {
		return "", Error.New("Failed to create chatter", err)
	}
	if err := app.ChatterChangeRoom(chatter.name, message.GetPayload()); err != nil {
		return "", Error.New("Failed to join room", err)
	}
	return Utilities.StringsToJsonObjectArray(app.GetRoomMessages(message.GetPayload())), nil
}

func (app *App) Leave(message *Message.Message) (string, error) {
	err := app.RemoveChatter(message.GetOrigin())
	if err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	return "", nil
}
