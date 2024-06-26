package appWebsocketHTTP

import (
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: app.PropagateMessage,
	}
}

func (app *AppWebsocketHTTP) PropagateMessage(client *Node.Node, message *Message.Message) error {
	client.WebsocketGroupcast(message.GetOrigin(), message)
	return nil
}
