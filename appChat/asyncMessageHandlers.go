package appChat

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *App) AddMessage(client *Client.Client, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[message.GetOrigin()]
	if chatter == nil {
		return Utilities.NewError("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Utilities.NewError("Room not found", nil)
	}
	chatMessage := NewChatMessage(chatter.id, message.GetPayload())
	room.AddMessage(chatMessage)
	client.AsyncMessage(topics.PROPAGATE_MESSAGE, chatter.roomId, chatMessage.Marshal())
	return nil
}
