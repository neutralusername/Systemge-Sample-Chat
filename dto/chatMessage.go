package dto

import (
	"encoding/json"
	"time"
)

const RINGBUFFER_SIZE = 7

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

func UnmarshalChatMessage(data string) *ChatMessage {
	var chatMessage ChatMessage
	json.Unmarshal([]byte(data), &chatMessage)
	return &chatMessage
}
