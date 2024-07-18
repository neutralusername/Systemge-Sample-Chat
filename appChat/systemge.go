package appChat

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Helpers"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Tcp"
	"SystemgeSampleChat/topics"
)

func (app *App) GetSystemgeComponentConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
	}
}

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
	return Helpers.StringsToJsonObjectArray(app.GetRoomMessages(message.GetPayload())), nil
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
