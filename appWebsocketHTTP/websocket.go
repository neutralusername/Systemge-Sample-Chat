package appWebsocketHTTP

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *AppWebsocketHTTP) AddMessage(node *Node.Node, connection *Node.WebsocketClient, message *Message.Message) error {
	err := node.AsyncMessage(topics.ADD_MESSAGE, dto.NewChatMessage(connection.GetId(), message.GetPayload()).Marshal())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to propagate message" + err.Error())
		}
	}
	return nil
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := node.AddToWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		websocketClient.Disconnect()
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to add to group" + err.Error())
		}
		return
	}
	responseChannel, err := node.SyncMessage(topics.JOIN, websocketClient.GetId())
	if err != nil {
		websocketClient.Disconnect()
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to join room" + err.Error())
		}
		return
	}
	response, err := responseChannel.ReceiveResponse()
	if err != nil {
		websocketClient.Disconnect()
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to receive response" + err.Error())
		}
		return
	}
	websocketClient.Send(Message.NewAsync("join", response.GetMessage().GetPayload()).Serialize())
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := node.RemoveFromWebsocketGroup("lobby", websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to remove from group" + err.Error())
		}
	}
	_, err = node.SyncMessage(topics.LEAVE, websocketClient.GetId())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log("Failed to leave room" + err.Error())
		}
	}
}
