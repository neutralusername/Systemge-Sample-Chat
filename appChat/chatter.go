package appChat

import "github.com/neutralusername/Systemge/Error"

type chatter struct {
	id     string //websocketId
	roomId string
}

func newChatter(id string) *chatter {
	return &chatter{
		id: id,
	}
}

func (app *App) addChatter(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.chatters[chatterId] != nil {
		return Error.New("Chatter already exists", nil)
	}
	app.chatters[chatterId] = newChatter(chatterId)
	return nil
}

func (app *App) removeChatter(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterId]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	if room := app.rooms[chatter.roomId]; room != nil {
		return Error.New("Chatter still in room", nil)
	}
	delete(app.chatters, chatterId)
	return nil
}
