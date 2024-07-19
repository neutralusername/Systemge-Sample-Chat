package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
	"SystemgeSampleChat/topics"
)

const LOGGER_PATH = "error.log"

func main() {
	Node.StartCommandLineInterface(true,
		Node.New(&Config.Node{
			Name: "nodeResolver",
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"Resolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"Resolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"Resolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"Resolver\"] ",
			},
		}, Resolver.New(&Config.Resolver{
			Server: &Config.TcpServer{
				Port:        60000,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ConfigServer: &Config.TcpServer{
				Port:        60001,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},

			TcpTimeoutMs: 5000,
		})),
		Node.New(&Config.Node{
			Name: "nodeBrokerChat",
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"Resolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"Resolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"Resolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"Resolver\"] ",
			},
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60002,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60003,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},

			SyncTopics:  []string{topics.LEAVE, topics.JOIN},
			AsyncTopics: []string{topics.ADD_MESSAGE},

			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name: "nodeBrokerWebsocketHTTP",
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"Resolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"Resolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"Resolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"Resolver\"] ",
			},
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60004,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60005,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},

			AsyncTopics: []string{topics.PROPAGATE_MESSAGE},

			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name: "nodeChat",
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"Resolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"Resolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"Resolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"Resolver\"] ",
			},
		}, appChat.New()),
		Node.New(&Config.Node{
			Name: "nodeWebsocketHTTP",
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"Resolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"Resolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"Resolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"Resolver\"] ",
			},
		}, appWebsocketHTTP.New()),
	)
}
