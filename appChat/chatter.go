package appChat

import "Systemge/Error"

type Chatter struct {
	id     string //websocketId
	roomId string
}

func NewChatter(id string) *Chatter {
	return &Chatter{
		id: id,
	}
}

func (app *App) AddChatter(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.chatters[chatterId] != nil {
		return Error.New("Chatter already exists", nil)
	}
	chatter := NewChatter(chatterId)
	app.chatters[chatterId] = chatter
	return nil
}

func (app *App) RemoveChatter(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterId]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room != nil {
		return Error.New("Chatter still in room", nil)
	}
	delete(app.chatters, chatterId)
	return nil
}
