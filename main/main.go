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
	//and if they are stopped first on stop, which is inherent to multi-module behaviour, the other modules will not be able to communicate with them during their stop/disconnect routine.
	//this demonstrates why multi modules are not always the best solution.
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs.
	Module.NewResolverServerFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerChat.systemge", ERROR_LOG_FILE_PATH).Start()
	Module.NewBrokerServerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()

	clientChat := Module.NewClient("clientApp", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appChat.New)
	clientWebsocket := Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		//order is important in this multi module because websocket app disconnects all clients when it stops and within onDisconnct() it communicates to the chat app that the client/chatter has disconnected.
		//if the chat app is stopped before the websocket app, the websocket app will not be able to communicate to the chat app that the client/chatter has disconnected.
		//which results in the chat app having chatters that will never be removed.
		clientWebsocket,
		clientChat,
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), Module.MergeCustomCommandHandlers(clientChat.GetApplication().GetCustomCommandHandlers(), clientWebsocket.GetWebsocketServer().GetCustomCommandHandlers()))
}
