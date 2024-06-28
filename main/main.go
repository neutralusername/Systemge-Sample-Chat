package main

import (
	"Systemge/Config"
	"Systemge/Module"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
)

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

	nodeChat := Module.NewNode(Config.Node{
		Name:       "nodeApp",
		LoggerPath: ERROR_LOG_FILE_PATH,
	}, appChat.New(), nil, nil)
	appWebsocketHTTP := appWebsocketHTTP.New()
	nodeWebsocket := Module.NewNode(Config.Node{
		Name:       "nodeWebsocketHTTP",
		LoggerPath: ERROR_LOG_FILE_PATH,
	}, appWebsocketHTTP, appWebsocketHTTP, appWebsocketHTTP)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		nodeWebsocket,
		nodeChat,
	))
}
