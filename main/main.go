package main

import (
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
	"SystemgeSampleChat/topics"

	"github.com/neutralusername/Systemge/Broker"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Resolver"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	loggerQueue := Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port: 8081,
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		NodeBrokerCounterIntervalMs:    1000,
		NodeResolverCounterIntervalMs:  1000,
		HeapUpdateIntervalMs:           1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.Node{
			Name:           "nodeResolver",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeResolver\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeResolver\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeResolver\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerChat\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerChat\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerChat\"] ", loggerQueue),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60002,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
				Domain:  "localhost",
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
				Domain:  "localhost",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60004,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
				Domain:  "localhost",
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
				Domain:  "localhost",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeChat",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeChat\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeChat\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeChat\"] ", loggerQueue),
		}, appChat.New()),
		Node.New(&Config.Node{
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeWebsocketHTTP\"] ", loggerQueue),
		}, appWebsocketHTTP.New()),
	),
	).StartBlocking()
}
