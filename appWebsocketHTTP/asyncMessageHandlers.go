package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: app.PropagateMessage,
	}
}

func (app *AppWebsocketHTTP) PropagateMessage(client *Client.Client, message *Message.Message) error {
	client.WebsocketGroupcast(message.GetOrigin(), message)
	return nil
}
