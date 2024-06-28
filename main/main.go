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
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs with their own command line interfaces.
	err := Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerChat.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Node.New(Config.Node{
			Name:       "nodeWebsocketHTTP",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, appWebsocketHTTP.New()),
		Node.New(Config.Node{
			Name:       "nodeApp",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, appChat.New()),
	))
}
