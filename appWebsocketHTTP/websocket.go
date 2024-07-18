package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() Config.Websocket {
	return Config.Websocket{
		Pattern: "/ws",
		Server: Config.TcpServer{
			Port: 8443,
		},
		HandleClientMessagesSequentially: false,

		ClientMessageCooldownMs: 0,
		ClientWatchdogTimeoutMs: 20000,
	}
}

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *AppWebsocketHTTP) AddMessage(node *Node.Node, connection *Node.WebsocketClient, message *Message.Message) error {
	err := node.AsyncMessage(topics.ADD_MESSAGE, connection.GetId(), message.GetPayload())
	if err != nil {
		node.GetLogger().Error(Error.New("Failed to send message", err).Error())
	}
	return nil
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := node.AddToWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		websocketClient.Disconnect()
		node.GetLogger().Error(Error.New("Failed to add to group", err).Error())
		return
	}
	response, err := node.SyncMessage(topics.JOIN, websocketClient.GetId(), "lobby")
	if err != nil {
		websocketClient.Disconnect()
		node.GetLogger().Error(Error.New("Failed to join room", err).Error())
		return
	}
	websocketClient.Send([]byte(response.Serialize()))
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := node.RemoveFromWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		node.GetLogger().Error(Error.New("Failed to remove from group", err).Error())
	}
	_, err = node.SyncMessage(topics.LEAVE, websocketClient.GetId(), "")
	if err != nil {
		node.GetLogger().Error(Error.New("Failed to leave room", err).Error())
	}
}
