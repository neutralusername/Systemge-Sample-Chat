package appChat

import "Systemge/Error"

type Chatter struct {
	name   string //websocketId
	roomId string
}

func NewChatter(name string) *Chatter {
	return &Chatter{
		name: name,
	}
}

func (app *App) AddChatter(chatterName string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.chatters[chatterName] != nil {
		return Error.New("Chatter already exists", nil)
	}
	chatter := NewChatter(chatterName)
	app.chatters[chatterName] = chatter
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
