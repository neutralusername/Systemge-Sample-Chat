package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Tcp"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
	"SystemgeSampleChat/topics"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Node.StartCommandLineInterface(true,
		Node.New(Config.Node{
			Name: "nodeResolver",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Resolver.New(Config.Resolver{
			Server:       Tcp.NewServer(60000, "MyCertificate.crt", "MyKey.key"),
			ConfigServer: Tcp.NewServer(60001, "MyCertificate.crt", "MyKey.key"),

			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerChat",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60002, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60002", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60003, "MyCertificate.crt", "MyKey.key"),

			SyncTopics:  []string{topics.LEAVE, topics.JOIN},
			AsyncTopics: []string{topics.ADD_MESSAGE},

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60004, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60004", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60005, "MyCertificate.crt", "MyKey.key"),

			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeChat",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appChat.New()),
		Node.New(Config.Node{
			Name: "nodeWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appWebsocketHTTP.New()),
	)
}
