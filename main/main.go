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
	//resolver and brokers are placed outside, because on startup they need to be started first.
	//and if they are closed first on stop, the other modules will not be able to communicate with them.
	//this demonstrates why multi modules are not always the best solution.
	Module.NewResolverServerFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerChat.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()

	clientChat := Module.NewClient("clientApp", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appChat.New)
	Module.StartCommandLineInterface(Module.NewMultiModule(

		//order is important in this multi module because websocket disconnects all clients when it stops and within the disconnect routine it communicates to the chat app that the client has disconnected.
		//if the chat app is stopped before the websocket client, the chat app will not be able to communicate to the websocket client that the client has disconnected.
		//which results in the websocket client having chatters that will never be removed.
		Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New),
		clientChat,
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), clientChat.GetApplication().GetCustomCommandHandlers())
}
