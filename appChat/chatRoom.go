package appChat

import (
	"SystemgeSampleChat/dto"

	"github.com/neutralusername/Systemge/Error"
)

type room struct {
	id string //websocketServer groupId

	//messages are stored in a ring buffer to limit memory usage per room
	messageRingBuffer           [dto.RINGBUFFER_SIZE]*dto.ChatMessage
	messageRingBufferWriteIndex int
	chatters                    map[string]*chatter //chatterId -> chatter
}

func rewRoom(id string) *room {
	return &room{
		id:                          id,
		messageRingBuffer:           [dto.RINGBUFFER_SIZE]*dto.ChatMessage{},
		messageRingBufferWriteIndex: 0,
		chatters:                    map[string]*chatter{},
	}
}

func (room *room) addMessage(message *dto.ChatMessage) {
	room.messageRingBuffer[room.messageRingBufferWriteIndex] = message
	room.messageRingBufferWriteIndex = (room.messageRingBufferWriteIndex + 1) % dto.RINGBUFFER_SIZE
}

func (app *App) getRoomMessages(roomId string) []string {
	room := app.rooms[roomId]
	if room == nil {
		return []string{}
	}
	messages := []string{}
	for i := 0; i < dto.RINGBUFFER_SIZE; i++ {
		if message := room.messageRingBuffer[(room.messageRingBufferWriteIndex+i)%dto.RINGBUFFER_SIZE]; message != nil {
			messages = append(messages, message.Marshal())
		}
	}
	return messages
}

func (app *App) addToRoom(chatterid string, roomId string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatter := app.chatters[chatterid]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[roomId]
	if room == nil {
		room = rewRoom(roomId)
		app.rooms[roomId] = room
	}
	chatter.roomId = roomId
	room.chatters[chatterid] = chatter
	return nil
}

func (app *App) removeFromRoom(chatterId string) error {
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
	chatter.roomId = ""
	return nil

}
