package main

import (
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHttp"
	"SystemgeSampleChat/topics"
	"time"

	"github.com/neutralusername/Systemge/BrokerResolver"
	"github.com/neutralusername/Systemge/BrokerServer"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
)

const LOGGER_PATH = "logs.log"

func main() {
	Helpers.StartPprof()
	if err := Dashboard.NewServer("DasbhboardServer",
		&Config.DashboardServer{
			HTTPServerConfig: &Config.HTTPServer{
				TcpServerConfig: &Config.TcpServer{
					Port: 8081,
				},
			},
			WebsocketServerConfig: &Config.WebsocketServer{
				Pattern:                 "/ws",
				ClientWatchdogTimeoutMs: 1000 * 60 * 10,
				TcpServerConfig: &Config.TcpServer{
					Port: 8444,
				},
			},
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60000,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			Metrics:                   true,
			Commands:                  true,
			SystemgeCommands:          true,
			HttpCommands:              true,
			WebsocketCommands:         true,
			HeapUpdateIntervalMs:      1000,
			GoroutineUpdateIntervalMs: 1000,
			StatusUpdateIntervalMs:    1000,
			MetricsUpdateIntervalMs:   1000,
			MaxChartEntries:           100,
		},
		nil, nil,
	).Start(); err != nil {
		panic(Error.New("Dashboard server failed to start", err))
	}
	if err := BrokerResolver.New("brokerResolver",
		&Config.MessageBrokerResolver{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60001,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
				ClientConfig: &Config.TcpClient{
					Address: "[::1]:60000",
				},
			},
			AsyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.PROPAGATE_MESSAGE: {
					Address: "[::1]:60002",
				},
				topics.ADD_MESSAGE: {
					Address: "[::1]:60002",
				},
			},
			SyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.JOIN: {
					Address: "[::1]:60002",
				},
				topics.LEAVE: {
					Address: "[::1]:60002",
				},
			},
		},
		nil, nil,
	).Start(); err != nil {
		panic(Error.New("MessageBroker resolver failed to start", err))
	}
	if err := BrokerServer.New("brokerServer",
		&Config.MessageBrokerServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60002,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE, topics.ADD_MESSAGE},
			SyncTopics:  []string{topics.JOIN, topics.LEAVE},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
				ClientConfig: &Config.TcpClient{
					Address: "[::1]:60000",
				},
			},
		},
		nil, nil,
	).Start(); err != nil {
		panic(Error.New("MessageBroker server failed to start", err))
	}
	appWebsocketHttp.New()
	appChat.New()
	<-make(chan time.Time)
}
