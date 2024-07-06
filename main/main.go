package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/TcpEndpoint"
	"Systemge/TcpServer"
	"Systemge/Utilities"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
	"SystemgeSampleChat/config"
	"SystemgeSampleChat/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	//resolver and brokers are placed outside, because on startup they need to be started first.
	//and if they are stopped first on stop, which is inherent to multi-module behaviour, the other modules will not be able to communicate with them during their stop/disconnect routine.
	//this demonstrates why multi modules are not always the best solution.
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs with their own command line interfaces.
	err := Resolver.New(Config.Resolver{
		Name:       config.RESOLVER_NAME,
		LoggerPath: ERROR_LOG_FILE_PATH,

		Server:   TcpServer.New(config.RESOLVER_PORT, config.CERT_PATH, config.KEY_PATH),
		Endpoint: TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),

		ConfigServer: TcpServer.New(config.RESOLVER_CONFIG_PORT, config.CERT_PATH, config.KEY_PATH),
	}).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.Broker{
		Name:                   "brokerChat",
		LoggerPath:             ERROR_LOG_FILE_PATH,
		DeliverImmediately:     true,
		ConnectionTimeoutMs:    3000,
		ResolverConfigEndpoint: TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_CONFIG_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
		SyncRequestTimeoutMs:   5000,

		Server:   TcpServer.New(config.BROKER_CHAT_PORT, config.CERT_PATH, config.KEY_PATH),
		Endpoint: TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.BROKER_CHAT_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),

		ConfigServer: TcpServer.New(config.BROKER_CHAT_CONFIG_PORT, config.CERT_PATH, config.KEY_PATH),

		SyncTopics:  []string{topics.JOIN, topics.LEAVE},
		AsyncTopics: []string{topics.ADD_MESSAGE},
	}).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.Broker{
		Name:                   "brokerWebsocketHTTP",
		LoggerPath:             ERROR_LOG_FILE_PATH,
		DeliverImmediately:     true,
		ConnectionTimeoutMs:    3000,
		ResolverConfigEndpoint: TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_CONFIG_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
		SyncRequestTimeoutMs:   5000,

		Server:   TcpServer.New(config.BROKER_WEBSOCKET_HTTP_PORT, config.CERT_PATH, config.KEY_PATH),
		Endpoint: TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.BROKER_WEBSOCKET_HTTP_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),

		ConfigServer: TcpServer.New(config.BROKER_WEBSOCKET_HTTP_CONFIG_PORT, config.CERT_PATH, config.KEY_PATH),

		SyncTopics:  []string{},
		AsyncTopics: []string{topics.PROPAGATE_MESSAGE},
	}).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Node.New(Config.Node{
			Name:                 "nodeWebsocketHTTP",
			LoggerPath:           ERROR_LOG_FILE_PATH,
			ResolverEndpoint:     TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncMessageTimeoutMs: 1000,
			HeartbeatIntervalMs:  100,
		}, appWebsocketHTTP.New()),
		Node.New(Config.Node{
			Name:                 "nodeApp",
			LoggerPath:           ERROR_LOG_FILE_PATH,
			ResolverEndpoint:     TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncMessageTimeoutMs: 1000,
			HeartbeatIntervalMs:  100,
		}, appChat.New()),
	))
}
