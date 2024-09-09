package appWebsocketHttp

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/BrokerClient"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	messageBrokerClient *BrokerClient.Client
	websocketServer     *WebsocketServer.WebsocketServer
	httpServer          *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}

	app.websocketServer = WebsocketServer.New("appWebsocketHttp",
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port: 8443,
			},
		},
		nil, nil,
		WebsocketServer.MessageHandlers{
			topics.ADD_MESSAGE: app.addMessage,
		},
		app.OnConnectHandler, app.OnDisconnectHandler,
	)
	app.httpServer = HTTPServer.New("appWebsocketHttp",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8080,
			},
		},
		nil, nil,
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

	commands := app.websocketServer.GetDefaultCommands()
	commands.Merge(app.httpServer.GetDefaultCommands())
	app.messageBrokerClient = BrokerClient.New("appWebsocketHttp",
		&Config.MessageBrokerClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{
				HeartbeatIntervalMs: 1000,
			},
			ResolverConnectionConfig: &Config.TcpSystemgeConnection{
				HeartbeatIntervalMs: 1000,
			},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
				ClientConfig: &Config.TcpClient{
					Address: "localhost:60000",
				},
			},
			ResolverClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
				},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},
		},
		messageHandler, commands,
	)

	if err := app.messageBrokerClient.Start(); err != nil {
		panic(err)
	}
	if err := app.websocketServer.Start(); err != nil {
		panic(err)
	}
	if err := app.httpServer.Start(); err != nil {
		panic(err)
	}

	return app
}

func (app *AppWebsocketHTTP) propagateMessage(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	responses := app.messageBrokerClient.SyncRequest(topics.JOIN, websocketClient.GetId())
	if len(responses) == 0 {
		return Error.New("Failed to receive response", nil)
	}
	response := responses[0]
	websocketClient.Send(Message.NewAsync("join", response.GetPayload()).Serialize())
	return nil
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketServer.WebsocketClient) {
	app.messageBrokerClient.SyncRequest(topics.LEAVE, websocketClient.GetId())
}

func (app *AppWebsocketHTTP) addMessage(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
	app.messageBrokerClient.AsyncMessage(topics.ADD_MESSAGE, dto.NewChatMessage(websocketClient.GetId(), message.GetPayload()).Marshal())
	return nil
}
