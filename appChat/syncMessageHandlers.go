package appChat

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.JOIN:  app.Join,
		topics.LEAVE: app.Leave,
	}
}

func (app *App) Join(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.AddChatter(message.GetOrigin()); err != nil {
		return "", Error.New("Failed to create chatter", err)
	}
	if err := app.AddToRoom(message.GetOrigin(), message.GetPayload()); err != nil {
		return "", Error.New("Failed to join room", err)
	}
	return Utilities.StringsToJsonObjectArray(app.GetRoomMessages(message.GetPayload())), nil
}

func (app *App) Leave(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.RemoveFromRoom(message.GetOrigin()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	if err := app.RemoveChatter(message.GetOrigin()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	return "", nil
}
