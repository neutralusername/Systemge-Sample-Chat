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
	"SystemgeSampleChat/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Module.StartCommandLineInterface(Module.NewMultiModule(true,
		Node.New(Config.Node{
			Name:   "nodeResolver",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Resolver.New(Config.Resolver{
			Server:       TcpServer.New(60000, "MyCertificate.crt", "MyKey.key"),
			ConfigServer: TcpServer.New(60001, "MyCertificate.crt", "MyKey.key"),

			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerChat",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60002, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60002", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
			ConfigServer: TcpServer.New(60003, "MyCertificate.crt", "MyKey.key"),

			SyncTopics:  []string{topics.LEAVE, topics.JOIN},
			AsyncTopics: []string{topics.ADD_MESSAGE},

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "example.com", Utilities.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeBrokerWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, Broker.New(Config.Broker{
			Server:       TcpServer.New(60004, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     TcpEndpoint.New("127.0.0.1:60004", "example.com", Utilities.GetFileContent("MyCertificate.crt")),
			ConfigServer: TcpServer.New(60005, "MyCertificate.crt", "MyKey.key"),

			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},

			ResolverConfigEndpoint: TcpEndpoint.New("127.0.0.1:60001", "example.com", Utilities.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name:   "nodeChat",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appChat.New()),
		Node.New(Config.Node{
			Name:   "nodeWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appWebsocketHTTP.New()),
	))
}
