package appChat

import (
	"Systemge/Application"
	"Systemge/Error"
)

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"chatters": app.GetChatters,
		"rooms":    app.GetRooms,
	}
}

func (app *App) GetChatters(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(args) != 1 {
		return Error.New("Invalid arguments", nil)
	}
	roomId := args[0]
	room := app.rooms[roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	for _, chatter := range room.chatters {
		println(chatter.id)
	}
	return nil
}

func (app *App) GetRooms(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for roomId := range app.rooms {
		println(roomId)
	}
	return nil
}
