package appChat

import (
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *App) AddMessage(node *Node.Node, message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	msg, err := Message.Deserialize([]byte(message.GetPayload()))
	if err != nil {
		return Error.New("Failed to deserialize message", err)
	}
	chatter := app.chatters[msg.GetTopic()]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	chatMessage := NewChatMessage(chatter.id, msg.GetPayload())
	room.AddMessage(chatMessage)
	propagateMsg := Message.NewAsync(chatter.roomId, chatMessage.Marshal())
	node.AsyncMessage(topics.PROPAGATE_MESSAGE, string(propagateMsg.Serialize()))
	return nil
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.JOIN:  app.Join,
		topics.LEAVE: app.Leave,
	}
}

func (app *App) Join(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.AddChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to create chatter", err)
	}
	if err := app.AddToRoom(message.GetPayload(), "lobby"); err != nil {
		return "", Error.New("Failed to join room", err)
	}
	return Helpers.StringsToJsonObjectArray(app.GetRoomMessages("lobby")), nil
}

func (app *App) Leave(node *Node.Node, message *Message.Message) (string, error) {
	if err := app.RemoveFromRoom(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	if err := app.RemoveChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	return "", nil
}
