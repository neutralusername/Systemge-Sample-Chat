package appWebsocketHttp

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/MessageBroker"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	messageBrokerClient *SystemgeConnection.SystemgeConnection
	websocketServer     *WebsocketServer.WebsocketServer
	httpServer          *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}

	app.websocketServer = WebsocketServer.New(
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpListenerConfig: &Config.TcpListener{
				Port: 8443,
			},
		},
		WebsocketServer.MessageHandlers{
			topics.ADD_MESSAGE: app.addMessage,
		},
		app.OnConnectHandler, app.OnDisconnectHandler,
	)
	app.httpServer = HTTPServer.New(
		&Config.HTTPServer{
			TcpListenerConfig: &Config.TcpListener{
				Port: 8080,
			},
		},
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)

	messageHandler := SystemgeConnection.NewConcurrentMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{
			topics.PROPAGATE_MESSAGE: app.propagateMessage,
		},
		SystemgeConnection.SyncMessageHandlers{},
		nil, nil,
	)

	messageBrokerClient, err := MessageBroker.NewMessageBrokerClient(
		&Config.MessageBrokerClient{
			Name:             "appWebsocketHttp",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60001",
			},
			DashboardClientConfig: &Config.DashboardClient{
				Name:             "appWebsocketHttp",
				ConnectionConfig: &Config.SystemgeConnection{},
				EndpointConfig: &Config.TcpEndpoint{
					Address: "localhost:60000",
				},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},
		},
		messageHandler, nil,
	)
	if err != nil {
		panic(err)
	}

	app.messageBrokerClient = messageBrokerClient

	if err := app.websocketServer.Start(); err != nil {
		panic(err)
	}
	if err := app.httpServer.Start(); err != nil {
		panic(err)
	}

	return app
}

func (app *AppWebsocketHTTP) propagateMessage(connection *SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	responseChannel, err := app.messageBrokerClient.SyncRequest(topics.JOIN, websocketClient.GetId())
	if err != nil {
		return Error.New("Failed to join room", err)
	}
	response := <-responseChannel
	if response == nil {
		return Error.New("Failed to receive response", err)
	}
	websocketClient.Send(Message.NewAsync("join", response.GetPayload()).Serialize())
	return nil
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketServer.WebsocketClient) {
	_, err := app.messageBrokerClient.SyncRequest(topics.LEAVE, websocketClient.GetId())
	if err != nil {
		panic(err)
	}
}

func (app *AppWebsocketHTTP) addMessage(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
	err := app.messageBrokerClient.AsyncMessage(topics.ADD_MESSAGE, dto.NewChatMessage(websocketClient.GetId(), message.GetPayload()).Marshal())
	if err != nil {
		panic(err)
	}
	return nil
}
