package appChat

import "Systemge/Error"

type Chatter struct {
	name   string //websocketId
	roomId string
}

func (app *App) AddChatter(chatterName string) (*Chatter, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.chatters[chatterName] != nil {
		return nil, Error.New("Chatter already exists", nil)
	}
	chatter := &Chatter{
		name: chatterName,
	}
	app.chatters[chatterName] = chatter
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
			//chat room is deleted once all chatters leave. that is why messages only persist while chatters are in the room
			delete(app.rooms, chatter.roomId)
		}
	}
	return nil
}
