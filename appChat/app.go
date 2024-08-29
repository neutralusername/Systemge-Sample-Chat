package appChat

import (
	"SystemgeSampleChat/topics"
	"sync"

	"github.com/neutralusername/Systemge/Commands"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/MessageBroker"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type App struct {
	rooms    map[string]*room    //roomId -> room
	chatters map[string]*chatter //chatterId -> chatter
	mutex    sync.Mutex

	messageBrokerClient *SystemgeConnection.SystemgeConnection
}

func New() *App {
	app := &App{
		mutex: sync.Mutex{},
	}

	messageBrokerClient, err := MessageBroker.NewMessageBrokerClient(
		&Config.MessageBrokerClient{
			Name:             "appChat",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60001",
			},
			DashboardClientConfig: &Config.DashboardClient{
				Name:             "appChat",
				ConnectionConfig: &Config.SystemgeConnection{},
				EndpointConfig: &Config.TcpEndpoint{
					Address: "localhost:60000",
				},
			},
			AsyncTopics: []string{topics.ADD_MESSAGE},
			SyncTopics:  []string{topics.JOIN, topics.LEAVE},
		},
		SystemgeConnection.NewConcurrentMessageHandler(
			SystemgeConnection.AsyncMessageHandlers{
				topics.ADD_MESSAGE: app.addMessage,
			},
			SystemgeConnection.SyncMessageHandlers{
				topics.JOIN:  app.join,
				topics.LEAVE: app.leave,
			},
			nil, nil,
		),
		Commands.Handlers{
			"getChatters": app.getChatters,
			"getRooms":    app.getRooms,
		},
	)
	if err != nil {
		panic(err)
	}
	app.messageBrokerClient = messageBrokerClient

	app.rooms = map[string]*room{}
	app.chatters = map[string]*chatter{}

	return app
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
