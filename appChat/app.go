package appChat

import (
	"sync"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	rooms    map[string]*room    //roomId -> room
	chatters map[string]*chatter //chatterId -> chatter
	mutex    sync.Mutex
}

func New() *App {
	app := &App{
		rooms:    map[string]*room{},
		chatters: map[string]*chatter{},
		mutex:    sync.Mutex{},
	}
	return app
}

func (app *App) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{
		"getChatters": app.getChatters,
		"getRooms":    app.getRooms,
	}
}

func (app *App) getChatters(node *Node.Node, args []string) (string, error) {
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

func (app *App) getRooms(node *Node.Node, args []string) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	resultStr := ""
	for roomId := range app.rooms {
		resultStr += roomId + ";"
	}
	return resultStr, nil
}
