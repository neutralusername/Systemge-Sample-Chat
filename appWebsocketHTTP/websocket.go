package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *AppWebsocketHTTP) AddMessage(client *Node.Node, connection *Node.WebsocketClient, message *Message.Message) error {
	err := client.AsyncMessage(topics.ADD_MESSAGE, connection.GetId(), message.GetPayload())
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to send message", err).Error())
	}
	return nil
}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
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

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	err := client.RemoveFromWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to remove from group", err).Error())
	}
	_, err = client.SyncMessage(topics.LEAVE, websocketClient.GetId(), "")
	if err != nil {
		client.GetLogger().Log(Error.New("Failed to leave room", err).Error())
	}
}
