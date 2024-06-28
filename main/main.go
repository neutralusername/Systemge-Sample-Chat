package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	//resolver and brokers are placed outside, because on startup they need to be started first.
	//and if they are stopped first on stop, which is inherent to multi-module behaviour, the other modules will not be able to communicate with them during their stop/disconnect routine.
	//this demonstrates why multi modules are not always the best solution.
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs.
	err := Resolver.New(Module.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Module.ParseBrokerConfigFromFile("brokerChat.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Module.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")).Start()
	if err != nil {
		panic(err)
	}
	nodeChat := Node.New(Config.Node{
		Name:       "nodeApp",
		LoggerPath: ERROR_LOG_FILE_PATH,
	}, appChat.New(), nil, nil)
	appWebsocketHTTP := appWebsocketHTTP.New()
	nodeWebsocket := Node.New(Config.Node{
		Name:       "nodeWebsocketHTTP",
		LoggerPath: ERROR_LOG_FILE_PATH,
	}, appWebsocketHTTP, appWebsocketHTTP, appWebsocketHTTP)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		nodeWebsocket,
		nodeChat,
	))
}
