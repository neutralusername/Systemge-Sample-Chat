package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topics.ADD_MESSAGE: func(connection *WebsocketClient.Client, message *Message.Message) error {
			err := app.client.AsyncMessage(message.GetTopic(), connection.GetId(), message.GetPayload())
			if err != nil {
				app.client.GetLogger().Log(Utilities.NewError("Failed to send message", err).Error())
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(connection *WebsocketClient.Client) {
	err := app.client.GetWebsocketServer().AddToGroup("lobby", connection.GetId())
	if err != nil {
		connection.Disconnect()
		app.client.GetLogger().Log(Utilities.NewError("Failed to add to group", err).Error())
	}
	response, err := app.client.SyncMessage(topics.JOIN, connection.GetId(), "lobby")
	if err != nil {
		connection.Disconnect()
		app.client.GetLogger().Log(Utilities.NewError("Failed to join room", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(connection *WebsocketClient.Client) {
	err := app.client.GetWebsocketServer().RemoveFromGroup("lobby", connection.GetId())
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Failed to remove from group", err).Error())
	}
	_, err = app.client.SyncMessage(topics.LEAVE, connection.GetId(), "")
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Failed to leave room", err).Error())
	}
}
