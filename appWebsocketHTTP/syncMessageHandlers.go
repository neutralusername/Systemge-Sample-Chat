package appWebsocketHTTP

import (
	"Systemge/Client"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{}
}
