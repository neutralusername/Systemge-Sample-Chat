package appChat

import (
	"Systemge/Error"
	"Systemge/Node"
)

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{
		"getChatters": app.GetChatters,
		"getRooms":    app.GetRooms,
	}
}

func (app *App) GetChatters(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(args) != 1 {
		return Error.New("Invalid arguments", nil)
	}
	room := app.rooms[args[0]]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	for _, chatter := range room.chatters {
		println(chatter.id)
	}
	return nil
}

func (app *App) GetRooms(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for roomId := range app.rooms {
		println(roomId)
	}
	return nil
}
