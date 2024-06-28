package appChat

import (
	"Systemge/Config"
	"Systemge/Node"
	"Systemge/Resolution"
	"Systemge/Utilities"
	"sync"
)

type App struct {
	rooms    map[string]*Room    //roomId -> room
	chatters map[string]*Chatter //chatterId -> chatter
	mutex    sync.Mutex
}

func New() *App {
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

func (app *App) GetApplicationConfig() Config.Application {
	return Config.Application{
		ResolverResolution:         Resolution.New("resolver", "127.0.0.1:60000", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
		HandleMessagesSequentially: false,
	}
}
