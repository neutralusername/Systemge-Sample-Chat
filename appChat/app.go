package appChat

import (
	"Systemge/Node"
	"sync"
)

type App struct {
	rooms    map[string]*Room    //roomId -> room
	chatters map[string]*Chatter //chatterId -> chatter
	mutex    sync.Mutex
}

func New() Node.Application {
	app := &App{
		rooms:    map[string]*Room{},
		chatters: map[string]*Chatter{},
		mutex:    sync.Mutex{},
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	//an alternative solution to the problem of async messages not being received by appChat during stoping using multi-modules would be to remove all remaining chatters and all rooms here
	return nil
}
