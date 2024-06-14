package appChat

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	Sender string    `json:"sender"`
	Text   string    `json:"text"`
	SentAt time.Time `json:"sentAt"`
}

func NewChatMessage(sender string, text string) *ChatMessage {
	return &ChatMessage{
		Sender: sender,
		Text:   text,
		SentAt: time.Now(),
	}
}

func (chatMessage *ChatMessage) Marshal() string {
	json, _ := json.Marshal(chatMessage)
	return string(json)
}

func (room *Room) AddMessage(message *ChatMessage) {
	room.messageRingBuffer[room.currentIndex] = message
	room.currentIndex = (room.currentIndex + 1) % RINGBUFFER_SIZE
}

func (app *App) GetRoomMessages(roomId string) []string {
	room := app.rooms[roomId]
	if room == nil {
		return []string{}
	}
	messages := []string{}
	for i := 0; i < RINGBUFFER_SIZE; i++ {
		index := (room.currentIndex + i) % RINGBUFFER_SIZE
		message := room.messageRingBuffer[index]
		if message != nil {
			messages = append(messages, message.Marshal())
		}
	}
	return messages
}
