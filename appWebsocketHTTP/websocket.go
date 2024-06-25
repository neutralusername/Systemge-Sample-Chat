package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Client.WebsocketMessageHandler {
	return map[string]Client.WebsocketMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *AppWebsocketHTTP) AddMessage(client *Client.Client, connection *Client.WebsocketClient, message *Message.Message) error {
	err := client.AsyncMessage(topics.ADD_MESSAGE, connection.GetId(), message.GetPayload())
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to send message", err).Error())
	}
	return nil
}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	err := client.AddToWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		websocketClient.Disconnect()
		client.GetLogger().Log(Error.New("Failed to add to group", err).Error())
	}
	response, err := client.SyncMessage(topics.JOIN, websocketClient.GetId(), "lobby")
	if err != nil {
		websocketClient.Disconnect()
		client.GetLogger().Log(Error.New("Failed to join room", err).Error())
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	err := client.RemoveFromWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to remove from group", err).Error())
	}
	_, err = client.SyncMessage(topics.LEAVE, websocketClient.GetId(), "")
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to leave room", err).Error())
	}
}
