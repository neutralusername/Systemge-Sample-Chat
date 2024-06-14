package appChat

import "Systemge/Error"

const RINGBUFFER_SIZE = 7

type Room struct {
	id string //websocketServer groupId

	//messages are stored in a ring buffer to limit memory usage per room
	messageRingBuffer           [RINGBUFFER_SIZE]*ChatMessage
	messageRingBufferWriteIndex int
	chatters                    map[string]*Chatter
}

func NewRoom(id string) *Room {
	return &Room{
		id:                          id,
		messageRingBuffer:           [RINGBUFFER_SIZE]*ChatMessage{},
		messageRingBufferWriteIndex: 0,
		chatters:                    map[string]*Chatter{},
	}
}

func (app *App) AddToRoom(chatterid string, roomId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterid]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[roomId]
	if room == nil {
		room = NewRoom(roomId)
		app.rooms[roomId] = room
	}
	chatter.roomId = roomId
	room.chatters[chatterid] = chatter
	return nil
}

func (app *App) RemoveFromRoom(chatterId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterId]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	delete(room.chatters, chatterId)
	if len(room.chatters) == 0 {
		delete(app.rooms, chatter.roomId)
	}
	return nil

}
