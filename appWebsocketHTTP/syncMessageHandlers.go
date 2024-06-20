package appWebsocketHTTP

import "Systemge/Application"

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}
