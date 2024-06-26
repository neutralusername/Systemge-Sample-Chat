package appChat

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChat/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *App) AddMessage(node *Node.Node, message *Message.Message) error {
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
	node.AsyncMessage(topics.PROPAGATE_MESSAGE, chatter.roomId, chatMessage.Marshal())
	return nil
}
