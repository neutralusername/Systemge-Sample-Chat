package appChat

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleChat/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *App) AddMessage(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	chatter := app.chatters[message.GetOrigin()]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	chatMessage := NewChatMessage(chatter.id, message.GetPayload())
	room.AddMessage(chatMessage)
	app.messageBrokerClient.AsyncMessage(topics.PROPAGATE_MESSAGE, chatter.roomId, chatMessage.Marshal())
	return nil
}
