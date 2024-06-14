package appChat

import (
	"Systemge/Application"
	"Systemge/MessageBrokerClient"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	logger              *Utilities.Logger
	messageBrokerClient *MessageBrokerClient.Client

	rooms    map[string]*Room    //roomId -> room
	chatters map[string]*Chatter //chatterId -> chatter
	mutex    sync.Mutex
}

func New(logger *Utilities.Logger, messageBrokerClient *MessageBrokerClient.Client) Application.Application {
	app := &App{
		logger:              logger,
		messageBrokerClient: messageBrokerClient,

		rooms:    map[string]*Room{},
		chatters: map[string]*Chatter{},
		mutex:    sync.Mutex{},
	}
	return app
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	//an alternative solution to the problem of async messages not being received by appChat during stoping using multi-modules would be to remove all remaining chatters and all rooms here
	return nil
}
