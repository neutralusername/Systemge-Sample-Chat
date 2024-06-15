package appChat

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.JOIN: app.Join,

		//leave needs to be sync because otherwise appChat will not receive the message using multi-modules as the onDisconnect() routine will only wait for the message to reach the broker with async messages and appChat will be stopped before the broker can route the message to the appChat module.
		//with sync messages, the onDisconnect() routine will wait for the response before finally stopping appWebsocket and afterwards appChat.
		topics.LEAVE: app.Leave,
	}
}

func (app *App) Join(message *Message.Message) (string, error) {
	if err := app.AddChatter(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to create chatter", err)
	}
	if err := app.AddToRoom(message.GetOrigin(), message.GetPayload()); err != nil {
		return "", Utilities.NewError("Failed to join room", err)
	}
	return Utilities.StringsToJsonObjectArray(app.GetRoomMessages(message.GetPayload())), nil
}

func (app *App) Leave(message *Message.Message) (string, error) {
	if err := app.RemoveFromRoom(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to leave room", err)
	}
	if err := app.RemoveChatter(message.GetOrigin()); err != nil {
		return "", Utilities.NewError("Failed to leave room", err)
	}
	return "", nil
}
