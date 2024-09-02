package appChat

import (
	"SystemgeSampleChat/dto"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

func (app *App) addMessage(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	chatMessage := dto.UnmarshalChatMessage(message.GetPayload())
	chatter := app.chatters[chatMessage.Sender]
	if chatter == nil {
		return
	}
	room := app.rooms[chatter.roomId]
	if room == nil {
		return
	}
	room.addMessage(chatMessage)
	app.messageBrokerClient.AsyncMessage(topics.PROPAGATE_MESSAGE, message.GetPayload())
}

func (app *App) join(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
	if err := app.addChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to create chatter", err)
	}
	if err := app.addToRoom(message.GetPayload(), "lobby"); err != nil {
		return "", Error.New("Failed to join room", err)
	}
	return Helpers.StringsToJsonObjectArray(app.getRoomMessages("lobby")), nil
}

func (app *App) leave(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
	if err := app.removeFromRoom(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	if err := app.removeChatter(message.GetPayload()); err != nil {
		return "", Error.New("Failed to leave room", err)
	}
	return "", nil
}
