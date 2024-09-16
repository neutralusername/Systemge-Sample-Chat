package appWebsocketHttp

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/BrokerClient"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
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
	if err := DashboardClientCustomService.New("appWebsocketHttp_websocketServer", &Config.DashboardClient{
		TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
			HeartbeatIntervalMs: 1000,
		},
		TcpClientConfig: &Config.TcpClient{
			Address: "[::1]:60000",
		},
	}, app.websocketServer, app.websocketServer.GetDefaultCommands()).Start(); err != nil {
		panic(Error.New("Dashboard client failed to start", err))
	}

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
	if err := DashboardClientCustomService.New("appWebsocketHttp_httpServer", &Config.DashboardClient{
		TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
			HeartbeatIntervalMs: 1000,
		},
		TcpClientConfig: &Config.TcpClient{
			Address: "[::1]:60000",
		},
	}, app.httpServer, app.httpServer.GetDefaultCommands()).Start(); err != nil {
		panic(Error.New("Dashboard client failed to start", err))
	}

	messageHandler := SystemgeMessageHandler.NewConcurrentMessageHandler(
		SystemgeMessageHandler.AsyncMessageHandlers{
			topics.PROPAGATE_MESSAGE: app.propagateMessage,
		},
		SystemgeMessageHandler.SyncMessageHandlers{},
		nil, nil,
	)

	commands := app.websocketServer.GetDefaultCommands()
	commands.Merge(app.httpServer.GetDefaultCommands())

	app.messageBrokerClient = BrokerClient.New("appWebsocketHttp",
		&Config.MessageBrokerClient{
			ResolutionAttemptRetryIntervalMs: 1000,
			ServerTcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
				HeartbeatIntervalMs: 1000,
			},
			ResolverTcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
				HeartbeatIntervalMs: 1000,
			},
			ResolverTcpClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
				},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},
		},
		messageHandler, commands,
	)
	if err := DashboardClientCustomService.New("appWebsocketHttp_brokerClient", &Config.DashboardClient{
		TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
			HeartbeatIntervalMs: 1000,
		},
		TcpClientConfig: &Config.TcpClient{
			Address: "[::1]:60000",
		},
	}, app.messageBrokerClient, app.messageBrokerClient.GetDefaultCommands()).Start(); err != nil {
		panic(Error.New("Dashboard client failed to start", err))
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
