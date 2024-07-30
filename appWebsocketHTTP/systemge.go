package appWebsocketHTTP

import (
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: app.PropagateMessage,
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *AppWebsocketHTTP) PropagateMessage(node *Node.Node, message *Message.Message) error {
	propagateMsg, err := Message.Deserialize([]byte(message.GetPayload()))
	if err != nil {
		return err
	}
	node.WebsocketGroupcast(propagateMsg.GetTopic(), Message.NewAsync("propagateMessage", propagateMsg.GetPayload()))
	return nil
}
