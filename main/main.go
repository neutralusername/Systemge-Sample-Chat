package main

import (
	"Systemge/Module"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = ":60000"
const HTTP_DEV_PORT = ":8080"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Module.NewResolverServerFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerChat.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()

	clientChat := Module.NewClient("clientApp", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appChat.New)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New),
		clientChat,
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), clientChat.GetApplication().GetCustomCommandHandlers())
}
