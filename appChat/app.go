package appChat

import (
	"SystemgeSampleChat/topics"
	"sync"

	"github.com/neutralusername/Systemge/Commands"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type App struct {
	rooms    map[string]*room    //roomId -> room
	chatters map[string]*chatter //chatterId -> chatter
	mutex    sync.Mutex

	systemgeClient *SystemgeClient.SystemgeClient
}

func New() *App {
	app := &App{
		mutex: sync.Mutex{},
	}
	messageHandler := SystemgeConnection.NewConcurrentMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{
			topics.ADD_MESSAGE: app.addMessage,
		},
		SystemgeConnection.SyncMessageHandlers{
			topics.JOIN:  app.join,
			topics.LEAVE: app.leave,
		},
		nil, nil,
	)

	app.systemgeClient = SystemgeClient.New(
		&Config.SystemgeClient{
			Name: "systemgeClient",
			EndpointConfigs: []*Config.TcpEndpoint{
				{
					Address: "localhost:60001",
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		func(connection *SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection *SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	Dashboard.NewClient(
		&Config.DashboardClient{
			Name:             "appChat",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60000",
			},
		},
		app.start, app.systemgeClient.Stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus,
		Commands.Handlers{
			"getChatters": app.getChatters,
			"getRooms":    app.getRooms,
		},
	)
	return app
}

func (app *App) start() error {
	err := app.systemgeClient.Start()
	if err != nil {
		return err
	}

	app.rooms = map[string]*room{}
	app.chatters = map[string]*chatter{}
	return nil
}

func (app *App) getChatters(args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(args) != 1 {
		return "", Error.New("Invalid arguments", nil)
	}
	room := app.rooms[args[0]]
	if room == nil {
		return "", Error.New("Room not found", nil)
	}
	resultStr := ""
	for _, chatter := range room.chatters {
		resultStr += chatter.id + ";"
	}
	return resultStr, nil
}

func (app *App) getRooms(args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	resultStr := ""
	for roomId := range app.rooms {
		resultStr += roomId + ";"
	}
	return resultStr, nil
}
