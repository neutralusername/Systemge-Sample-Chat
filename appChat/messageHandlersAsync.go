package appChat

import (
	"Systemge/Application"
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleChat/topics"
	"time"
)

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.ADD_MESSAGE: app.AddMessage,
	}
}

func (app *App) AddMessage(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	chatter := app.chatters[message.GetOrigin()]
	if chatter == nil {
		return Error.New("Chatter not found", nil)
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return Error.New("Room not found", nil)
	}
	chatMessage := &ChatMessage{
		Name:   chatter.name,
		SentAt: time.Now(),
		Text:   message.GetPayload(),
	}
	room.messageRingBuffer[room.currentIndex] = chatMessage
	room.currentIndex = (room.currentIndex + 1) % RINGBUFFER_SIZE
	app.messageBrokerClient.AsyncMessage(topics.PROPAGATE_MESSAGE, chatter.roomId, chatMessage.Marshal())
	return nil
}
