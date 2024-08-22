package appWebsocketHttp

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
	"github.com/neutralusername/Systemge/SystemgeServer"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	systemgeServer  *SystemgeServer.SystemgeServer
	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}
	app.systemgeServer = SystemgeServer.New(&Config.SystemgeServer{
		Name: "systemgeServer",
		ListenerConfig: &Config.SystemgeListener{
			TcpListenerConfig: &Config.TcpListener{
				Port: 60001,
			},
		},
		ConnectionConfig: &Config.SystemgeConnection{},
	}, nil, nil,
		SystemgeMessageHandler.New(SystemgeMessageHandler.AsyncMessageHandlers{
			topics.PROPAGATE_MESSAGE: app.propagateMessage,
		}, SystemgeMessageHandler.SyncMessageHandlers{}))
	app.websocketServer = WebsocketServer.New(&Config.WebsocketServer{
		ClientWatchdogTimeoutMs: 1000 * 60,
		Pattern:                 "/ws",
		TcpListenerConfig: &Config.TcpListener{
			Port: 8443,
		},
	}, WebsocketServer.MessageHandlers{
		topics.ADD_MESSAGE: app.addMessage,
	}, app.OnConnectHandler, app.OnDisconnectHandler)
	app.httpServer = HTTPServer.New(&Config.HTTPServer{
		TcpListenerConfig: &Config.TcpListener{
			Port: 8080,
		},
	}, HTTPServer.Handlers{
		"/": HTTPServer.SendDirectory("../frontend"),
	})
	Dashboard.NewClient(&Config.DashboardClient{
		Name:             "appWebsocketHttp",
		ConnectionConfig: &Config.SystemgeConnection{},
		EndpointConfig: &Config.TcpEndpoint{
			Address: "localhost:60000",
		},
	}, app.start, app.stop, app.systemgeServer.GetMetrics, app.getStatus, nil)
	return app
}

func (app *AppWebsocketHTTP) getStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.systemgeServer.Start(); err != nil {
		return Error.New("Failed to start systemgeServer", err)
	}
	if err := app.websocketServer.Start(); err != nil {
		app.systemgeServer.Stop()
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.systemgeServer.Stop()
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.systemgeServer.Stop()
	app.status = Status.STOPPED
	return nil
}

func (app *AppWebsocketHTTP) propagateMessage(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) WebsocketPropagate(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	responseChannel, err := app.systemgeServer.SyncRequest(topics.JOIN, websocketClient.GetId())
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
	_, err := app.systemgeServer.SyncRequest(topics.LEAVE, websocketClient.GetId())
	if err != nil {
		panic(err)
	}
}

func (app *AppWebsocketHTTP) addMessage(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
	err := app.systemgeServer.AsyncMessage(topics.ADD_MESSAGE, dto.NewChatMessage(websocketClient.GetId(), message.GetPayload()).Marshal())
	if err != nil {
		panic(err)
	}
	return nil
}
