package appChat

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.addMessage,
	}
}

func (app *App) addMessage(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatMessage := dto.UnmarshalChatMessage(message.GetPayload())
	chatter := app.chatters[chatMessage.Sender]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	room.addMessage(chatMessage)
	node.AsyncMessage(topics.PROPAGATE_MESSAGE, message.GetPayload())
	return nil
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.JOIN:  app.join,
		topics.LEAVE: app.leave,
	}
}

func (app *App) join(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.addChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to create chatter", err)
	}
	if err := app.addToRoom(message.GetPayload(), "lobby"); err != nil {
		return "", Error.New("Failed to join room", err)
	}
	return Helpers.StringsToJsonObjectArray(app.getRoomMessages("lobby")), nil
}

func (app *App) leave(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.removeFromRoom(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	if err := app.removeChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	return "", nil
}
