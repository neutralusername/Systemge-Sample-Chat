package appChat

import (
	"Systemge/Application"
	"Systemge/Utilities"
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
		return Utilities.NewError("Invalid arguments", nil)
	}
	room := app.rooms[args[0]]
	if room == nil {
		return Utilities.NewError("Room not found", nil)
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
