package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/TcpEndpoint"
	"Systemge/Utilities"
	"SystemgeSampleChat/topics"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: TcpEndpoint.New("127.0.0.1:60000", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_MESSAGE: app.PropagateMessage,
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *AppWebsocketHTTP) PropagateMessage(node *Node.Node, message *Message.Message) error {
	node.WebsocketGroupcast(message.GetOrigin(), message)
	return nil
}
