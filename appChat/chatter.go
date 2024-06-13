package appChat

import "Systemge/Error"

type Chatter struct {
	name   string
	roomId string
}

func (app *App) AddChatter(chatterId string) (*Chatter, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.chatters[chatterId] != nil {
		return nil, Error.New("Chatter already exists", nil)
	}
	chatter := &Chatter{
		name: chatterId,
	}
	app.chatters[chatterId] = chatter
	return chatter, nil
}

func (app *App) RemoveChatter(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterId]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	delete(app.chatters, chatterId)
	room := app.rooms[chatter.roomId]
	if room != nil {
		delete(room.chatters, chatterId)
		if len(room.chatters) == 0 {
			delete(app.rooms, chatter.roomId)
		}
	}
	return nil
}
