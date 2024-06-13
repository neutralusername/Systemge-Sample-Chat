package appChat

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	Name   string    `json:"name"`
	Text   string    `json:"text"`
	SentAt time.Time `json:"sentAt"`
}

func NewChatMessage(name string, text string) *ChatMessage {
	return &ChatMessage{
		Name:   name,
		Text:   text,
		SentAt: time.Now(),
	}
}

func (chatMessage *ChatMessage) Marshal() string {
	json, _ := json.Marshal(chatMessage)
	return string(json)
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
