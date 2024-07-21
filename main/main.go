package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Dashboard"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Tools"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
	"SystemgeSampleChat/topics"
)

const LOGGER_PATH = "logs.log"

func main() {
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
		ErrorLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Error \"dashboard\"] ",
		},
		WarningLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Warning \"dashboard\"] ",
		},
		InfoLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Info \"dashboard\"] ",
		},
		DebugLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Debug \"dashboard\"] ",
		},
	}, Dashboard.New(&Config.Dashboard{
		Http: &Config.Http{
			Server: &Config.TcpServer{
				Port:        8081,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
		},
		Pattern:                "/",
		StatusUpdateIntervalMs: 1000,
	},
		Node.New(&Config.Node{
			Name:           "nodeResolver",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeResolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeResolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeResolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeResolver\"] ",
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
			Name:           "nodeBrokerChat",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeBrokerChat\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeBrokerChat\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeBrokerChat\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeBrokerChat\"] ",
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
			Name:           "nodeBrokerWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeBrokerWebsocketHTTP\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeBrokerWebsocketHTTP\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeBrokerWebsocketHTTP\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeBrokerWebsocketHTTP\"] ",
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
			Name:           "nodeChat",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeChat\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeChat\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeChat\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeChat\"] ",
			},
		}, appChat.New()),
		Node.New(&Config.Node{
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeWebsocketHTTP\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeWebsocketHTTP\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeWebsocketHTTP\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeWebsocketHTTP\"] ",
			},
		}, appWebsocketHTTP.New()),
	),
	).Start()
	<-make(chan struct{})
}
