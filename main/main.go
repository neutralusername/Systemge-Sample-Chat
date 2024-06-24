package main

import (
	"Systemge/Module"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	//resolver and brokers are placed outside, because on startup they need to be started first.
	//and if they are stopped first on stop, which is inherent to multi-module behaviour, the other modules will not be able to communicate with them during their stop/disconnect routine.
	//this demonstrates why multi modules are not always the best solution.
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs.
	err := Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerChat.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	clientChat := Module.NewClient(&Module.ClientConfig{
		Name:                   "clientApp",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCertPath:    RESOLVER_TLS_CERT_PATH,
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appChat.New, nil)
	clientWebsocket := Module.NewCompositeClientWebsocketHTTP(&Module.ClientConfig{
		Name:                   "clientWebsocketHTTP",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCertPath:    RESOLVER_TLS_CERT_PATH,
		WebsocketPattern:       "/ws",
		WebsocketPort:          WEBSOCKET_PORT,
		HTTPPort:               HTTP_PORT,
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appWebsocketHTTP.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		//order is important in this multi module because websocket app disconnects all clients when it stops and within onDisconnct() it communicates to the chat app that the client/chatter has disconnected.
		//if the chat app is stopped before the websocket app, the websocket app will not be able to communicate to the chat app that the client/chatter has disconnected.
		//which results in the chat app having chatters that will never be removed.
		clientWebsocket,
		clientChat,
	), clientChat.GetApplication().GetCustomCommandHandlers(), clientWebsocket.GetWebsocketServer().GetCustomCommandHandlers())
}
