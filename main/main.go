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
		Name:                   config.BROKER_CHAT_NAME,
		LoggerPath:             ERROR_LOG_FILE_PATH,
		DeliverImmediately:     true,
		NodeTimeoutMs:          3000,
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
		Name:                   config.BROKER_WEBSOCKET_HTTP_NAME,
		LoggerPath:             ERROR_LOG_FILE_PATH,
		DeliverImmediately:     true,
		NodeTimeoutMs:          3000,
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
			Name:                      config.NODE_WEBSOCKET_HTTP_NAME,
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			BrokerHeartbeatIntervalMs: 100,
		}, appWebsocketHTTP.New()),
		Node.New(Config.Node{
			Name:                      config.NODE_CHAT_NAME,
			LoggerPath:                ERROR_LOG_FILE_PATH,
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_ADDRESS+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			BrokerHeartbeatIntervalMs: 100,
		}, appChat.New()),
	))
}
