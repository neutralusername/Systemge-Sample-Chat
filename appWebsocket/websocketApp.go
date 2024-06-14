package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleChat/topics"
)

type WebsocketApp struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) Application.WebsocketApplication {
	return &WebsocketApp{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,
	}
}

func (app *WebsocketApp) OnStart() error {
	return nil
}

func (app *WebsocketApp) OnStop() error {
	return nil
}

func (app *WebsocketApp) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: func(message *Message.Message) error {
			app.messageBrokerClient.GetWebsocketServer().Groupcast(message.GetOrigin(), message)
			return nil
		},
	}
}

func (app *WebsocketApp) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}

func (app *WebsocketApp) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *WebsocketApp) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		topics.ADD_MESSAGE: func(connection *WebsocketClient.Client, message *Message.Message) error {
			err := app.messageBrokerClient.AsyncMessage(message.GetTopic(), connection.GetId(), message.GetPayload())
			if err != nil {
				app.logger.Log(Error.New("Failed to send message", err).Error())
			}
			return nil
		},
	}
}

func (app *WebsocketApp) OnConnectHandler(connection *WebsocketClient.Client) {
	err := app.messageBrokerClient.GetWebsocketServer().AddToGroup("lobby", connection.GetId())
	if err != nil {
		connection.Disconnect()
		app.logger.Log(Error.New("Failed to add to group", err).Error())
	}
	response, err := app.messageBrokerClient.SyncMessage(topics.JOIN, connection.GetId(), "lobby")
	if err != nil {
		connection.Disconnect()
		app.logger.Log(Error.New("Failed to join room", err).Error())
	}
	connection.Send([]byte(response.Serialize()))
}

func (app *WebsocketApp) OnDisconnectHandler(connection *WebsocketClient.Client) {
	err := app.messageBrokerClient.GetWebsocketServer().RemoveFromGroup("lobby", connection.GetId())
	if err != nil {
		app.logger.Log(Error.New("Failed to remove from group", err).Error())
	}
	_, err = app.messageBrokerClient.SyncMessage(topics.LEAVE, connection.GetId(), "")
	if err != nil {
		app.logger.Log(Error.New("Failed to leave room", err).Error())
	}
}
