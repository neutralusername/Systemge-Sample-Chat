package appChat

import "Systemge/Error"

const RINGBUFFER_SIZE = 7

type Room struct {
	id string //websocketServer groupId

	//messages are stored in a ring buffer to limit memory usage per room
	messageRingBuffer [RINGBUFFER_SIZE]*ChatMessage
	currentIndex      int
	chatters          map[string]*Chatter
}

func NewRoom(id string) *Room {
	return &Room{
		id:                id,
		messageRingBuffer: [RINGBUFFER_SIZE]*ChatMessage{},
		currentIndex:      0,
		chatters:          map[string]*Chatter{},
	}
}

func (app *App) AddToRoom(chatterName string, roomId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterName]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[roomId]
	if room == nil {
		room = NewRoom(roomId)
		app.rooms[roomId] = room
	}
	chatter.roomId = roomId
	room.chatters[chatterName] = chatter
	return nil
}

func (app *App) RemoveFromRoom(chatterName string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterName]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room != nil {
		delete(room.chatters, chatterName)
		if len(room.chatters) == 0 {
			delete(app.rooms, chatter.roomId)
		}
	}
	return nil

}
