package appChat

import (
	"Systemge/Application"
	"Systemge/Client"
	"sync"
)

type App struct {
	client *Client.Client

	rooms    map[string]*Room    //roomId -> room
	chatters map[string]*Chatter //chatterId -> chatter
	mutex    sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client: client,

		rooms:    map[string]*Room{},
		chatters: map[string]*Chatter{},
		mutex:    sync.Mutex{},
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	//an alternative solution to the problem of async messages not being received by appChat during stoping using multi-modules would be to remove all remaining chatters and all rooms here
	return nil
}
