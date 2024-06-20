package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: func(message *Message.Message) error {
			app.client.GetWebsocketServer().Groupcast(message.GetOrigin(), message)
			return nil
		},
	}
}
