package appWebsocketHTTP

import (
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: app.propagateMessage,
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *AppWebsocketHTTP) propagateMessage(node *Node.Node, message *Message.Message) error {
	node.WebsocketBroadcast(message)
	return nil
}
